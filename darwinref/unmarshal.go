// Unmarshal the Darwin Reference XML
package darwinref

import (
  bolt "github.com/coreos/bbolt"
  "encoding/xml"
  "log"
  "time"
)

func (r *DarwinReference) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {
  return r.internalUpdate( func( tx *bolt.Tx ) error {
    return r.unmarshalXML( tx, decoder, start )
  })
}

func (r *DarwinReference) unmarshalXML( tx *bolt.Tx, decoder *xml.Decoder, start xml.StartElement ) error {
  r.toc = make( map[string]*Toc )
  crs := r.newCrsImport()
  tplcount := 0
  r.lateRunningReasons = make( map[int]string )
  r.cancellationReasons = make( map[int]string )
  r.cisSource = make( map[string]string )
  r.via = make( map[string][]*Via )

  for _, attr := range start.Attr {
    switch attr.Name.Local {
    case "timetableId":
      r.timetableId = attr.Value
    }
  }

  // Reason map to write to
  var late bool
  var inReason bool

  for {
    token, err := decoder.Token()
    if err != nil {
      return err
    }

    switch tok := token.(type) {
    case xml.StartElement:
      switch tok.Name.Local {
      case "LocationRef":
        var loc *Location = &Location{}
        if err = decoder.DecodeElement( loc, &tok ); err != nil {
          return err
        }
        loc.Date = time.Now()

        if err, updated := r.addTiploc( loc ); err != nil {
          return err
        } else if updated {
          tplcount ++
        }

        // Append to CRS map
        crs.append( loc )

      case "TocRef":
        var toc *Toc = &Toc{}
        if err = decoder.DecodeElement( toc, &tok ); err != nil {
          return err
        }

        if _, exists := r.toc[ toc.Toc ]; exists {
          log.Println( "Toc", toc.Toc, "duplicated" )
        } else {
          r.toc[ toc.Toc ] = toc
        }

      case "LateRunningReasons":
        inReason = true
        late = true

      case "CancellationReasons":
        inReason = true
        late = false

      case "Reason":
        if inReason {
          var reason *Reason = &Reason{}
          if err = decoder.DecodeElement( reason, &tok ); err != nil {
            return err
          }
          if late {
            r.lateRunningReasons[ reason.Code ] = reason.Text
          } else {
            r.cancellationReasons[ reason.Code ] = reason.Text
          }
        }

      case "CISSource":
        var cis *CISSource = &CISSource{}
        if err = decoder.DecodeElement( cis, &tok ); err != nil {
          return err
        }
        r.cisSource[ cis.Code ] = cis.Name

      case "Via":
        var via *Via = &Via{}
        if err = decoder.DecodeElement( via, &tok ); err != nil {
          return err
        }

        var key string = via.At + "," + via.Dest
        if  arr, exists := r.via[ key ]; exists {
          // To complicate things there are some duplicate entries
          exists = false
          for _, ent := range arr {
            exists = exists || via.Equals( ent )
          }
          if !exists {
            r.via[ key ] = append( arr, via )
          }
        } else {
          r.via[ key ] = append( r.via[ key ], via )
        }

      default:
        log.Println( "Unknown element", tok.Name.Local )
      }

    case xml.EndElement:
      if !inReason {
        log.Printf( "Imported %d Tiplocs", tplcount )

        if err, count := crs.write(); err != nil {
          return err
        } else {
          log.Printf( "Imported %d CRS", count )
        }

        return nil
      }
      inReason = false
    }
  }

}
