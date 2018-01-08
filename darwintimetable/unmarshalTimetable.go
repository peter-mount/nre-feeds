// Reference timetable
package darwintimetable

import (
  "encoding/xml"
  "log"
  //"time"
)

func (t *Timetable) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {
  t.Journeys = make( map[string]*Journey )
  var assocs []*Association

  for _, attr := range start.Attr {
    switch attr.Name.Local {
    case "timetableID":
      t.TimetableId = attr.Value
    }
  }

  for {
    token, err := decoder.Token()
    if err != nil {
      return err
    }

    switch tok := token.(type) {
    case xml.StartElement:

      switch tok.Name.Local {
      case "Journey":
        var j *Journey = &Journey{}

        if err = decoder.DecodeElement( j, &tok ); err != nil {
          return err
        }

        if _, ok := t.Journeys[ j.RID ]; ok {
          log.Println( "RID", j.RID, "duplicated" )
        } else {
          t.Journeys[ j.RID ] = j
        }

      case "Association":
        var a *Association = &Association{}

        if err = decoder.DecodeElement( a, &tok ); err != nil {
          return err
        }

        assocs = append( assocs, a )

      default:
        log.Println( "Unknown element", tok.Name.Local )
      }

    case xml.EndElement:

      for _, a := range assocs {
        if j1, ok := t.Journeys[a.Main.RID]; ok {
          if j2, ok := t.Journeys[a.Assoc.RID]; ok {
            j1.Associations = append( j1.Associations, a )
            j2.Associations = append( j2.Associations, a )
            } else {
              log.Println( "Assoc", a.Assoc.RID, "not found, main", a.Main.RID )
            }
            } else {
              log.Println( "Main", a.Main.RID, "not found, Assoc", a.Assoc.RID )
            }
          }

      return nil
    }
  }
}
