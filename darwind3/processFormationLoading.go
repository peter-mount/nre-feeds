package darwind3

import (
	"log"
)

// Process processes an inbound loading element containing train formation data.
func (l *Loading) Process(tx *Transaction) error {

	// Retrieve the schedule to be updated
	sched := tx.d3.GetSchedule(l.RID)

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

	sched.Date = tx.pport.TS
	if tx.d3.PutSchedule(sched) {
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
			tx.d3.updateAssociations(sched)
			return
		}
	}

	log.Println("Unknown Loading", l)
}
