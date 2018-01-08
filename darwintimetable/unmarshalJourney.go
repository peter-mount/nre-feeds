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
      j.SSD = attr.Value
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
      var elementStruct interface{}
      switch tok.Name.Local {
      case "OR":
        elementStruct = &OR{}
      case "OPOR":
        elementStruct = &OPOR{}
      case "IP":
        elementStruct = &IP{}
      case "OPIP":
        elementStruct = &OPIP{}
      case "PP":
        elementStruct = &PP{}
      case "DT":
        elementStruct = &DT{}
      case "OPDT":
        elementStruct = &OPDT{}
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

        elementStruct = nil
      default:
        log.Println( "Unknown element", tok.Name.Local, j.RID, j.SSD )
        elementStruct = nil
      }

      if elementStruct != nil {
        if err = decoder.DecodeElement( elementStruct, &tok ); err != nil {
          return err
        }

        j.Schedule = append( j.Schedule, elementStruct )
      }

    case xml.EndElement:
      return nil
    }
  }

}
