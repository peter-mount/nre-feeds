package darwind3

import (
  "encoding/xml"
  "log"
)

// Notification that a Train Schedule is now deactivated in Darwin.
type DeactivatedSchedule struct {
  XMLName     xml.Name  `json:"-" xml:"deactivated"`
  RID         string    `xml:"rid,attr"`
}

// Processor interface
func (p *DeactivatedSchedule) Process( d3 *DarwinD3, r *Pport ) error {
  log.Printf( "Deactivated rid=%s\n", p.RID )

  /*
  if p.UR != nil {
    return p.UR.Process( d3, p )
  }
  */
  return nil
}
