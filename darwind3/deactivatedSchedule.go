package darwind3

import (
  "encoding/xml"
)

// Notification that a Train Schedule is now deactivated in Darwin.
type DeactivatedSchedule struct {
  XMLName     xml.Name  `json:"-" xml:"deactivated"`
  RID         string    `xml:"rid,attr"`
}

// Processor interface
func (p *DeactivatedSchedule) Process( tx *Transaction ) error {
  tx.DeleteSchedule( p.RID )
  return nil
}
