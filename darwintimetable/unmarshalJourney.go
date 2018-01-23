// Reference timetable
package darwintimetable

import (
  "encoding/xml"
  "log"
  "strconv"
  //"time"
)

func (j *Journey) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {
  for _, attr := range start.Attr {
    switch attr.Name.Local {
    case "rid":
      j.RID = attr.Value
    case "uid":
      j.UID = attr.Value
    case "trainId":
      j.TrainID = attr.Value
    case "ssd":
      j.SSD.Parse( attr.Value )
    case "toc":
      j.Toc = attr.Value
    case "trainCat":
      j.TrainCat = attr.Value
    case "isPassengerSvc":
      if b, err := strconv.ParseBool( attr.Value ); err != nil {
        return err
      } else {
        j.Passenger = b
      }
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
      case "OR":
        if err = decodeAndAppendLocation( decoder, &tok, j, &OR{} ); err != nil {
          return err
        }

      case "OPOR":
        if err = decodeAndAppendLocation( decoder, &tok, j, &OPOR{} ); err != nil {
          return err
        }

      case "IP":
        if err = decodeAndAppendLocation( decoder, &tok, j, &IP{} ); err != nil {
          return err
        }

      case "OPIP":
        if err = decodeAndAppendLocation( decoder, &tok, j, &OPIP{} ); err != nil {
          return err
        }

      case "PP":
        if err = decodeAndAppendLocation( decoder, &tok, j, &PP{} ); err != nil {
          return err
        }

      case "DT":
        if err = decodeAndAppendLocation( decoder, &tok, j, &DT{} ); err != nil {
          return err
        }

      case "OPDT":
        if err = decodeAndAppendLocation( decoder, &tok, j, &OPDT{} ); err != nil {
          return err
        }

      case "cancelReason":
        var cr *cancelReason = &cancelReason{}
        if err = decoder.DecodeElement( cr, &tok ); err != nil {
          return err
        }
        if cr.text != "" {
          if crv, err := strconv.Atoi( cr.text ); err != nil {
            return err
            } else {
              j.CancelReason = crv
            }
        }

      default:
        log.Println( "Unknown element", tok.Name.Local, j.RID, j.SSD )
      }

    case xml.EndElement:
      return nil
    }
  }

}

func decodeAndAppendLocation( decoder *xml.Decoder, tok *xml.StartElement, j *Journey, v interface{ Location() *Location } ) error {
  if err := decoder.DecodeElement( v, tok ); err != nil {
    return err
  }
  j.Schedule = append( j.Schedule, v.Location() )
  return nil
}
