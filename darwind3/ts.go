package darwind3

import (
  "encoding/xml"
  "log"
)

// TS Train Status
type TS struct {
  XMLName           xml.Name  `json:"-" xml:"TS"`
  RID               string    `json:"rid" xml:"rid,attr"`
  UID               string    `json:"uid" xml:"uid,attr"`
  SSD               string    `json:"ssd" xml:"ssd,attr"`
  ReverseFormation  bool      `json:"isReverseFormation,omitempty" xml:"isReverseFormation,attr,omitempty"`

}

// Processor interface
func (p *TS) Process( d3 *DarwinD3, r *Pport ) error {
  log.Printf(
    "TS rid=%s uid=%s ssd=%s reverse=%s\n",
    p.RID,
    p.UID,
    p.SSD,
    p.ReverseFormation )

  /*
  if p.UR != nil {
    return p.UR.Process( d3, p )
  }
  */
  return nil
}
