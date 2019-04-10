package darwind3

import "github.com/etcd-io/bbolt"

// Process processes an inbound StationMessage
func (sm *StationMessage) Process(tx *Transaction) error {
	if tx.d3.cache.tx != nil {
		return sm.process(tx, tx.d3.cache.tx)
	}
	return tx.d3.Update(func(dbtx *bbolt.Tx) error {
		return sm.process(tx, dbtx)
	})
}

func (sm *StationMessage) process(tx *Transaction, dbtx *bbolt.Tx) error {
	sm.Date = tx.pport.TS

	old := tx.d3.Messages.get(dbtx, sm.ID)
	_ = tx.d3.Messages.put(dbtx, sm)

	tx.d3.EventManager.PostEvent(&DarwinEvent{
		Type:                   Event_StationMessage,
		ExistingStationMessage: old,
		NewStationMessage:      sm,
	})

	return nil
}
