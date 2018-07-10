package darwinkb

import (
  "encoding/json"
  "github.com/peter-mount/golib/kernel/bolt"
  "log"
  "errors"
)

const (
  stationXml = "station.xml"
  stationJson = "station.json"
)
func (r *DarwinKB) refreshStations() error {

  updateRequired, err := r.refreshFile( stationXml, "https://datafeeds.nationalrail.co.uk/api/staticfeeds/4.0/stations" )
  if err != nil {
    return err
  }

  // If no update check to see if the bucket is empty forcing an update
  if !updateRequired {
    err = r.boltDb.View( func( tx *bolt.Tx ) error {
      bucket := tx.Bucket( "stations" )
      if bucket == nil {
        return errors.New( "Bucket not found" )
      }
      cursor := bucket.Cursor()
      k, _ := cursor.First()
      if k == "" {
        updateRequired = true
      }
      return nil
    })
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

  a := make( map[string]map[string]interface{} )
  err = json.Unmarshal( b.Bytes(), &a )
  if err != nil {
    return err
  }

  stations := a["StationList"]["Station"].([]interface{})
  log.Printf( "Found %d stations", len(stations) )

  err = r.boltDb.Update( func( tx *bolt.Tx ) error {
    bucket := tx.Bucket( "stations" )
    if bucket == nil {
      return errors.New( "Bucket not found" )
    }

    // Clear existing entries
    log.Println( "Clearing existing entries" )
    err = bucket.ForEach( func(k string, v []byte) error {
      return bucket.Delete( k )
    })
    if err != nil {
      return err
    }

    // Insert entries one per entry using CrsCode as the key
    for _, c := range stations {

      d := c.(map[string]interface{})
      crs := d["CrsCode"].(string)

      b, err := json.Marshal( d )
      if err != nil {
        return err
      }

      err = bucket.Put( crs, b )
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
