// Unmarshal the Darwin Reference XML
package darwinref

import (
  "encoding/xml"
  "log"
  //"strconv"
  //"time"
)

func (r *DarwinReference) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {
  r.Tiploc = make( map[string]*Location )
  r.Crs = make( map[string][]*Location )
  r.Toc = make( map[string]*Toc )
  r.LateRunningReasons = make( map[int]string )
  r.CancellationReasons = make( map[int]string )
  r.CISSource = make( map[string]string )
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

        if _, exists := r.Tiploc[ loc.Tiploc ]; exists {
          log.Println( "Tiploc", loc.Tiploc, "duplicated" )
        } else {
          r.Tiploc[ loc.Tiploc ] = loc
        }

        if loc.Crs != "" {
          r.Crs[ loc.Crs ] = append( r.Crs[ loc.Crs ], loc )
        }

      case "TocRef":
        var toc *Toc = &Toc{}
        if err = decoder.DecodeElement( toc, &tok ); err != nil {
          return err
        }

        if _, exists := r.Toc[ toc.Toc ]; exists {
          log.Println( "Toc", toc.Toc, "duplicated" )
        } else {
          r.Toc[ toc.Toc ] = toc
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
          r.LateRunningReasons[ reason.Code ] = reason.Text
          } else {
            r.CancellationReasons[ reason.Code ] = reason.Text
          }
        }

      case "CISSource":
        var cis *CISSource = &CISSource{}
        if err = decoder.DecodeElement( cis, &tok ); err != nil {
          return err
        }
        r.CISSource[ cis.Code ] = cis.Name

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
        return nil
      }
      inReason = false
    }
  }

}
