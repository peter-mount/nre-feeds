package darwinkb

import (
  "github.com/peter-mount/golib/kernel/bolt"
  "log"
)

const (
  incidentXml = "incident.xml"
  incidentJson = "incident.json"
)

// An entry in the incident index
type IncidentEntry struct {
  Id      string  `json:"id"`
  Summary string  `json:"summary"`
}

func (r *DarwinKB) GetIncidents() ([]byte, error) {
  // Works as we have the index as a single key
  b, err := r.GetIncident( "index" )
  return b, err
}

func (r *DarwinKB) GetIncident( id string ) ([]byte, error) {
  var data []byte
  err := r.View( "incidents", func( bucket *bolt.Bucket ) error {
    data = bucket.Get( id )
    return nil
  } )
  return data, err
}

func (r *DarwinKB) refreshIncidents() {
  err := r.refreshIncidentsImpl()
  if err != nil {
    log.Println( "refreshIncidents:", err )
  }
}

func (r *DarwinKB) refreshIncidentsImpl() error {

  updateRequired, err := r.refreshFile( incidentXml, "https://datafeeds.nationalrail.co.uk/api/staticfeeds/5.0/incidents" )
  if err != nil {
    return err
  }

  // If no update check to see if the bucket is empty forcing an update
  if !updateRequired {
    updateRequired, err = r.bucketEmpty( "incidents" )
    if err != nil {
      return err
    }
  }

  // Give up if no update is required
  if !updateRequired {
    return nil
  }

  b, err := r.xml2json( incidentXml, incidentJson )
  if err != nil {
    return err
  }
  log.Println( "Parsing JSON" )

  root, err := unmarshalBytes( b )
  if err != nil {
    return err
  }

  incidents, _ := GetJsonArray( root, "Incidents", "PtIncident" )
  log.Println( "Found", len(incidents), "incidents" )

  err = r.Update( "incidents", func( bucket *bolt.Bucket ) error {
    err := bucketRemoveAll( bucket )
    if err != nil {
      return err
    }

    // slice containing index of all entries
    var index []*IncidentEntry

    for _, incident := range incidents {
      o := incident.(map[string]interface{})
      indexEntry := &IncidentEntry{
        Id: o["IncidentNumber"].(string),
        Summary: o[ "Summary" ].(string),
      }
      index = append( index, indexEntry )

      // The individual entry
      err = bucket.PutJSON( indexEntry.Id, incident )
      if err != nil {
        return err
      }
    }

    // Now the index entry
    err = bucket.PutJSON( "index", index )
    if err != nil {
      return err
    }

    return nil
  } )
  if err != nil {
    return err
  }

  log.Printf( "Updated %d incidents", len(incidents) )
  return nil
}
