package darwind3

import (
	"encoding/xml"
	"github.com/etcd-io/bbolt"
)

// Notification that a Train schedule is now deactivated in Darwin.
type DeactivatedSchedule struct {
	XMLName xml.Name `json:"-" xml:"deactivated"`
	RID     string   `xml:"rid,attr"`
}

// Processor interface
func (p *DeactivatedSchedule) Process(tx *Transaction) error {
	return tx.d3.UpdateBulkAware(func(dbtx *bbolt.Tx) error {
		return p.process(tx, dbtx)
	})
}

func (p *DeactivatedSchedule) process(tx *Transaction, dbtx *bbolt.Tx) error {

	// Get the affected schedule
	sched := GetSchedule(dbtx, p.RID)

	// Delete it if we have one
	if sched != nil {
		// Mark as not active & persist
		sched.Active = false
		sched.Date = tx.pport.TS
		PutSchedule(dbtx, sched)

		tx.d3.updateAssociations(dbtx, sched)
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
