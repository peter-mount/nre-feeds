package darwind3

// Process inbound associations
func (a *Association) Process(tx *Transaction) error {
	a.Date = tx.pport.TS
	a.Assoc.Times.UpdateTime()
	a.Main.Times.UpdateTime()

	err := a.Main.processSched(tx, a)
	if err != nil {
		return err
	}

	return a.Assoc.processSched(tx, a)
}

func (as *AssocService) processSched(tx *Transaction, a *Association) error {

	assocs := tx.d3.GetAssociations(as.RID)
	if assocs == nil {
		assocs = &Associations{RID: as.RID}
	}

	found := false
	for i, assoc := range assocs.Associations {
		if assoc.Equals(a) {
			assocs.Associations[i] = a
			found = true
		}
	}
	if !found {
		assocs.Associations = append(assocs.Associations, a)
	}
	assocs.putAssociations(tx)

	sched := tx.d3.GetSchedule(assocs.RID)
	if sched != nil {

		tx.d3.UpdateAssociations(sched)

		sched.Date = tx.pport.TS
		tx.d3.EventManager.PostEvent(&DarwinEvent{
			Type:     Event_ScheduleUpdated,
			RID:      sched.RID,
			Schedule: sched,
		})
	}

	return nil
}
