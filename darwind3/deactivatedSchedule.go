package darwind3

import (
	"encoding/xml"
)

// Notification that a Train Schedule is now deactivated in Darwin.
type DeactivatedSchedule struct {
	XMLName xml.Name `json:"-" xml:"deactivated"`
	RID     string   `xml:"rid,attr"`
}

// Processor interface
func (p *DeactivatedSchedule) Process(tx *Transaction) error {

	// Get the affected schedule
	sched := tx.d3.GetSchedule(p.RID)

	// Delete it if we have one
	if sched != nil {
		// Mark as not active & persist
		sched.Active = false
		sched.Date = tx.pport.TS
		tx.d3.putSchedule(sched)
	}

	// Post event
	tx.d3.EventManager.PostEvent(&DarwinEvent{
		Type: Event_Deactivated,
		RID:  p.RID,
		// This is ok if nil but helps listeners know what to remove
		Schedule: sched,
	})

	return nil
}
