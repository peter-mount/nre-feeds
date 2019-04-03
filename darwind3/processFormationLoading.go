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

	sched.Update(func() error {

		for _, loc := range sched.Locations {
			if loc.Tiploc == l.Tiploc && l.Times.Equals(&loc.Times) {
				loc.Loading = l
				tx.d3.updateAssociations(sched)
				return nil
			}
		}

		log.Println("Unknown Loading", l)
		return nil
	})

	sched.Date = tx.pport.TS
	tx.d3.putSchedule(sched)
	tx.d3.EventManager.PostEvent(&DarwinEvent{
		Type:     Event_ScheduleUpdated,
		RID:      sched.RID,
		Schedule: sched,
	})
	return nil
}
