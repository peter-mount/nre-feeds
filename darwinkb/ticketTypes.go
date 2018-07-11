package darwinkb

import (
  "github.com/peter-mount/golib/kernel/bolt"
  "github.com/peter-mount/sortfold"
  "log"
  "sort"
  "time"
)

const (
  ticketTypesXml = "ticket-types.xml"
  ticketTypesJson = "ticket-types.json"
)

type TicketType struct {
  Id    string  `json:"id"`
  Name  string  `json:"name"`
}

func (r *DarwinKB) GetTicketTypes() ([]byte, error) {
  // Works as we have the index as a single key
  b, err := r.GetTicketType( "index" )
  return b, err
}

func (r *DarwinKB) GetTicketIDs() ([]byte, error) {
  // Works as we have the index as a single key
  b, err := r.GetTicketType( "idIndex" )
  return b, err
}

func (r *DarwinKB) GetTicketType( id string ) ([]byte, error) {
  var data []byte
  err := r.View( "ticketTypes", func( bucket *bolt.Bucket ) error {
    data = bucket.Get( id )
    return nil
  } )
  return data, err
}

func (r *DarwinKB) refreshTicketTypes() {
  err := r.refreshTicketTypesImpl()
  if err != nil {
    log.Println( "refreshTicketTypes:", err )
  }
}

func (r *DarwinKB) refreshTicketTypesImpl() error {

  updateRequired, err := r.refreshFile( ticketTypesXml, "https://datafeeds.nationalrail.co.uk/api/staticfeeds/4.0/ticket-types", 6 * time.Hour)
  if err != nil {
    return err
  }

  // If no update check to see if the bucket is empty forcing an update
  if !updateRequired {
    updateRequired, err = r.bucketEmpty( "ticketTypes" )
    if err != nil {
      return err
    }
  }

  // Give up if no update is required
  if !updateRequired {
    return nil
  }

  b, err := r.xml2json( ticketTypesXml, ticketTypesJson )
  if err != nil {
    return err
  }

  log.Println( "Parsing JSON" )

  root, err := unmarshalBytes( b )
  if err != nil {
    return err
  }

  ForceJsonArray( root, "TicketTypeDescriptionList", "TicketTypeDescription", "ApplicableTocs", "IncludedTocs", "TocRef" )

  var index []*TicketType
  var idIndex []*TicketType

  ticketTypes, _ := GetJsonArray( root, "TicketTypeDescriptionList", "TicketTypeDescription" )
  log.Println( "Found", len(ticketTypes), "ticketTypes" )

  err = r.Update( "ticketTypes", func( bucket *bolt.Bucket ) error {
    err := bucketRemoveAll( bucket )
    if err != nil {
      return err
    }

    for _, ticket := range ticketTypes {
      o := ticket.(map[string]interface{})

      code, _ := GetJsonObjectValue( o, "TicketTypeCode" )
      name, _ := GetJsonObjectValue( o, "TicketTypeName" )
      id, _ := GetJsonObjectValue( o, "TicketTypeIdentifier" )

      err = bucket.PutJSON( id.(string), ticket )
      if err != nil {
        return err
      }
      idIndex = append( idIndex, &TicketType{ Id: id.(string), Name: name.(string) } )

      if s, ok := code.(string); ok {
        err = bucket.PutJSON( s, ticket )
        if err != nil {
          return err
        }
        index = append( index, &TicketType{ Id: s, Name: name.(string) } )
      } else if a, ok := code.([]interface{}); ok {
        for _, c := range a {
          err = bucket.PutJSON( c.(string), ticket )
          if err != nil {
            return err
          }
          index = append( index, &TicketType{ Id: c.(string), Name: name.(string) } )
        }
      } else {
        log.Println( "Unknown:", code )
      }

    }

    sort.SliceStable( index, func( i, j int ) bool { return sortfold.CompareFold( index[i].Name, index[j].Name ) < 0 } )

    err = bucket.PutJSON( "index", index )
    if err != nil {
      return err
    }

    err = bucket.PutJSON( "idIndex", idIndex )
    if err != nil {
      return err
    }
    return nil
  } )
  if err != nil {
    return err
  }

  log.Printf( "Updated %d ticketTypes", len(index) + len(idIndex) )
  return nil
}
