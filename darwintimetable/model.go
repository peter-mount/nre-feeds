// Reference timetable
package darwintimetable

import (
  "encoding/xml"
  "log"
  "strconv"
  //"time"
)

type PportTimetable struct {
  TimetableId     string      `xml:"timetableID,attr"`
  Journeys      []*Journey    `xml:"Journey"`
}

type Journey struct {
  RID             string      `xml:"rid,attr"`
  UID             string      `xml:"uid,attr"`
  TrainID         string      `xml:"trainId"`
  SSD             string      `xml:"ssd,attr"`
  Toc             string      `xml:"toc,attr"`
  TrainCat        string      `xml:"trainCat,attr"`
  Passenger       bool        `xml:"isPassengerSvc,attr"`
  // The schedule
  Schedule      []interface{} `xml:,any`
  CancelReason    int         `xml:"cancelReason"`
}

type OR struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // CallPtAttributes
  Pta       string      `xml:"pta,attr"`
  Ptd       string      `xml:"ptd,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // False destination to be used at this location
  FalseDest string      `xml:"fd,attr"`
}

type OPOR struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
}

type IP struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // CallPtAttributes
  Pta       string      `xml:"pta,attr"`
  Ptd       string      `xml:"ptd,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
  // False destination to be used at this location
  FalseDest string      `xml:"fd,attr"`
}

type OPIP struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

type PP struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wtp       string      `xml:"wtp,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

type DT struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // CallPtAttributes
  Pta       string      `xml:"pta,attr"`
  Ptd       string      `xml:"ptd,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

type OPDT struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

type cancelReason struct {
  text string `xml:",chardata"`
}

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
