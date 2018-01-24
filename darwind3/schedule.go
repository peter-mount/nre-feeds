package darwind3

import (
  "encoding/xml"
  "log"
)

// Train Schedule
type Schedule struct {
  RID               string
  UID               string
  TrainId           string
  SSD               string
  Toc               string
  // Default P
  Status            string
  // Default OO
  TrainCat          string
  // Default true
  PassengerService  bool
  // Default true
  Active            bool
  // Default false
  Deleted           bool
  // Default false
  Charter           bool
  // Cancel reason
  CancelReason      DisruptionReason
  // The locations in this schedule
  Locations       []*Location
}

func (s *Schedule) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {
  // Defaults
  s.Status = "P"
  s.TrainCat = "OO"
  s.PassengerService = true
  s.Active = true

  for _, attr := range start.Attr {
    switch attr.Name.Local {
      case "rid":
        s.RID = attr.Value

      case "uid":
        s.UID = attr.Value

      case "trainId":
        s.TrainId = attr.Value

      case "ssd":
        s.SSD = attr.Value

      case "toc":
        s.Toc = attr.Value

      case "status":
        s.Status = attr.Value

      case "isPassengerSvc":
        s.PassengerService = attr.Value == "true"

      case "isActive":
        s.Active = attr.Value == "true"

      case "deleted":
        s.Deleted = attr.Value == "true"

      case "isCharter":
        s.Charter = attr.Value == "true"
    }
  }

  for {
    token, err := decoder.Token()
    if err != nil {
      return err
    }

    switch tok := token.(type) {
      case xml.StartElement:
        var elem *Location
        switch tok.Name.Local {
          case "OR":
            elem = &Location{ Type: "OR" }

          case "OPOR":
            elem = &Location{ Type: "OPOR" }

          case "IP":
            elem = &Location{ Type: "IP" }

          case "OPIP":
            elem = &Location{ Type: "OPIP" }

          case "PP":
            elem = &Location{ Type: "PP" }

          case "DT":
            elem = &Location{ Type: "DT" }

          case "OPDT":
            elem = &Location{ Type: "OPDT" }

          case "cancelReason":
            if err := decoder.DecodeElement( &s.CancelReason, &tok ); err != nil {
              return err
            }

          default:
            if err := decoder.Skip(); err != nil {
              return err
            }
        }

        if elem != nil {
          if err := decoder.DecodeElement( elem, &tok ); err != nil {
            return err
          }
          s.Locations = append( s.Locations, elem )
        }

      case xml.EndElement:
        return nil
    }
  }
}

// Processor interface
func (p *Schedule) Process( d3 *DarwinD3, r *Pport ) error {
  log.Printf(
    "Schedule rid=%s uid=%s trainId=%s ssd=%s toc=%s status=%s trainCat=%s isPax=%v active=%v deleted=%v charter=%v cancelReason=%v locs=%d\n",
    p.RID,
    p.UID,
    p.TrainId,
    p.SSD,
    p.Toc,
    p.Status,
    p.TrainCat,
    p.PassengerService,
    p.Active,
    p.Deleted,
    p.Charter,
    p.CancelReason,
    len( p.Locations ) )

  /*
  if p.UR != nil {
    return p.UR.Process( d3, p )
  }
  */
  return nil
}
