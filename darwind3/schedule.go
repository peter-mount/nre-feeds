package darwind3

import (
  "darwintimetable"
  "encoding/xml"
  "fmt"
  "log"
)

// Train Schedule
type Schedule struct {
  XMLName     xml.Name  `json:"-" xml:"schedule"`
  RID         string    `xml:"rid,attr"`
  UID         string    `xml:"uid,attr"`
  TrainId     string    `xml:"trainId,attr"`
  SSD         string    `xml:"ssd,attr"`
  Toc         string    `xml:"toc,attr"`
  // Default P
  Status      string    `xml:"status,attr,omitempty"`
  // Default OO
  TrainCat          string    `xml:"trainCat,attr,omitempty"`
  // Default true
  PassengerService  bool      `xml:"isPassengerSvc,attr,omitempty"`
  // Default true
  Active            bool      `xml:"isActive,attr,omitempty"`
  // Default false
  Deleted           bool      `xml:"deleted,attr,omitempty"`
  // Default false
  Charter           bool      `xml:"isCharter,attr,omitempty"`
  // Minimum of 2 entries
  Schedule       []*darwintimetable.Location
}

// Processor interface
func (p *Schedule) Process( d3 *DarwinD3, r *Pport ) error {
  log.Printf(
    "Schedule rid=%s uid=%s trainId=%s ssd=%s toc=%s status=%s trainCat=%s isPax=%s active=%s deleted=%s charter=%s\n",
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
    p.Charter )

  /*
  if p.UR != nil {
    return p.UR.Process( d3, p )
  }
  */
  return nil
}

func (p *Schedule) String() string {
  s := fmt.Sprintf(
    "Schedule rid=%s uid=%s trainId=%s ssd=%s toc=%s status=%s trainCat=%s isPax=%s active=%s deleted=%s charter=%s\n",
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
    p.Charter )

   return s
}
