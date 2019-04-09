package darwind3

// Process inbound associations
func (a *Association) Process(tx *Transaction) error {
	a.Date = tx.pport.TS

	err := a.Main.processSched(tx, a)
	if err != nil {
		return err
	}

	return a.Assoc.processSched(tx, a)
}

func (as *AssocService) processSched(tx *Transaction, a *Association) error {
	sched := tx.d3.GetSchedule(as.RID)

	// No schedule then try to fetch it from the timetable
	if sched == nil {
		sched = tx.ResolveSchedule(as.RID)
	}

	// Still no schedule then We've got a TS for a train with no known schedule so create one
	if sched == nil {
		sched = &Schedule{RID: as.RID}
		sched.Defaults()
	}

	sched.appendAssociation(a)

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

func (sched *Schedule) appendAssociation(a *Association) {
	// Replace if we already have it
	for i, e := range sched.Associations {
		if e.Equals(a) {
			sched.Associations[i] = a
			return
		}
	}

	// Not found then add it
	sched.Associations = append(sched.Associations, a)
}
