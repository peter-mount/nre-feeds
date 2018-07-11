package darwinkb

import (
  "github.com/peter-mount/golib/kernel/bolt"
  "log"
)

const (
  stationXml = "station.xml"
  stationJson = "station.json"
)

func (r *DarwinKB) GetStation( crs string ) ([]byte, error) {
  var data []byte
  err := r.View( "stations", func( bucket *bolt.Bucket ) error {
    data = bucket.Get( crs )
    return nil
  } )
  return data, err
}

func (r *DarwinKB) refreshStations() {
  err := r.refreshStationsImpl()
  if err != nil {
    log.Println( "refreshStations:", err )
  }
}

func (r *DarwinKB) refreshStationsImpl() error {

  updateRequired, err := r.refreshFile( stationXml, "https://datafeeds.nationalrail.co.uk/api/staticfeeds/4.0/stations" )
  if err != nil {
    return err
  }

  // If no update check to see if the bucket is empty forcing an update
  if !updateRequired {
    updateRequired, err = r.bucketEmpty( "stations" )
    if err != nil {
      return err
    }
  }

  // Give up if no update is required
  if !updateRequired {
    return nil
  }

  b, err := r.xml2json( stationXml, stationJson )
  if err != nil {
    return err
  }

  log.Println( "Parsing JSON" )

  root, err := unmarshalBytes( b )
  if err != nil {
    return err
  }

  stations, _ := GetJsonArray( root, "StationList", "Station" )
  log.Printf( "Found %d stations", len(stations) )

  err = r.Update( "stations", func( bucket *bolt.Bucket ) error {
    err := bucketRemoveAll( bucket )
    if err != nil {
      return err
    }

    // Insert entries one per entry using CrsCode as the key
    for _, c := range stations {

      d := c.(map[string]interface{})
      crs := d["CrsCode"].(string)

      err = bucket.PutJSON( crs, d )
      if err != nil {
        return err
      }
    }

    return nil
  } )
  if err != nil {
    return err
  }

  log.Printf( "Updated %d stations", len(stations) )
  return nil
}
