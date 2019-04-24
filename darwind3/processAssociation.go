package darwind3

import (
	"github.com/etcd-io/bbolt"
)

// Process inbound associations
func (a *Association) Process(tx *Transaction) error {
	return tx.d3.UpdateBulkAware(func(dbtx *bbolt.Tx) error {
		return a.process(tx, dbtx)
	})
}

func (a *Association) process(tx *Transaction, dbtx *bbolt.Tx) error {
	a.Date = tx.pport.TS
	a.Assoc.Times.UpdateTime()
	a.Main.Times.UpdateTime()

	err := a.Main.processSched(tx, dbtx, a)
	if err != nil {
		return err
	}

	return a.Assoc.processSched(tx, dbtx, a)
}

func (as *AssocService) processSched(tx *Transaction, dbtx *bbolt.Tx, a *Association) error {

	assocs := getAssociations(dbtx, as.RID)
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
	assocs.putAssociations(dbtx)

	sched := GetSchedule(dbtx, assocs.RID)
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
