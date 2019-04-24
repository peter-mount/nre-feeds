package darwind3

import (
	"github.com/etcd-io/bbolt"
	"log"
)

// Process processes an inbound loading element containing train formation data.
func (l *Loading) Process(tx *Transaction) error {
	return tx.d3.UpdateBulkAware(func(dbtx *bbolt.Tx) error {
		return l.process(tx, dbtx)
	})
}

func (l *Loading) process(tx *Transaction, dbtx *bbolt.Tx) error {
	// TODO remove this once we have received loading
	log.Println("Loading received!", l.RID)

	// Retrieve the schedule to be updated
	sched := GetSchedule(dbtx, l.RID)

	// No schedule then try to fetch it from the timetable
	if sched == nil {
		sched = tx.ResolveSchedule(l.RID)
	}

	// If no schedule then warn as we need UID & SSD but don't have it in the
	// Loading message
	if sched == nil {
		log.Println("Unknown RID in Loading", l.RID, "fid", l.Fid)
		return nil
	}

	// Set the Date field to the TS time
	l.Date = tx.pport.TS

	sched.appendFormationLoading(tx, l)
	tx.d3.updateAssociations(dbtx, sched)

	sched.Date = tx.pport.TS
	if PutSchedule(dbtx, sched) {
		tx.d3.EventManager.PostEvent(&DarwinEvent{
			Type:     Event_ScheduleUpdated,
			RID:      sched.RID,
			Schedule: sched,
		})
	}
	return nil
}

func (sched *Schedule) appendFormationLoading(tx *Transaction, l *Loading) {
	for _, loc := range sched.Locations {
		if loc.Tiploc == l.Tiploc && l.Times.Equals(&loc.Times) {
			loc.Loading = l
			return
		}
	}

	log.Println("Unknown Loading", l)
}
