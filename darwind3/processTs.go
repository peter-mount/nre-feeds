package darwind3

// Process processes an inbound Train Status update, merging it with an existing
// schedule in the database
func (p *TS) Process(tx *Transaction) error {

	// Retrieve the schedule to be updated
	sched := tx.d3.GetSchedule(p.RID)

	// No schedule then try to fetch it from the timetable
	if sched == nil {
		sched = tx.ResolveSchedule(p.RID)
	}

	// Still no schedule then We've got a TS for a train with no known schedule so create one
	if sched == nil {
		sched = &Schedule{
			RID: p.RID,
			UID: p.UID,
			SSD: p.SSD,
		}
		sched.Defaults()
	}

	// If the TS is older than what's in the schedule then then do nothing as it's
	// presumably old data that's sent out of sync or it's from a snapshot
	if tx.pport.TS.Before(sched.Date) {
		return nil
	}

	// set forecast date of the new entries
	for _, a := range p.Locations {
		a.Forecast.Date = tx.pport.TS
	}

	// SnapshotUpdate the LateReason
	sched.LateReason = p.LateReason

	// Run through schedule locations, any that match the new ones update the forecast
	for _, a := range sched.Locations {
		for _, b := range p.Locations {
			if a.EqualInSchedule(b) {
				a.MergeFrom(b)
			}
		}
	}

	// Append any locations not in the schedule
	sortRequired := false
	for _, a := range p.Locations {
		f := true
		for _, b := range sched.Locations {
			if a.EqualInSchedule(b) {
				f = false
			}
		}
		if f {
			sched.Locations = append(sched.Locations, a)
			sortRequired = true
		}
	}

	tx.d3.updateAssociations(sched)

	// Sort if required else just update the times
	if sortRequired {
		sched.Sort()
	} else {
		sched.UpdateTime()
	}

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
