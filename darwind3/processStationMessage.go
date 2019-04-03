package darwind3

// Process processes an inbound StationMessage
func (sm *StationMessage) Process(tx *Transaction) error {
	sm.Date = tx.pport.TS

	old := tx.d3.Messages.Get(sm.ID)
	tx.d3.Messages.Put(sm)

	tx.d3.EventManager.PostEvent(&DarwinEvent{
		Type:                   Event_StationMessage,
		ExistingStationMessage: old,
		NewStationMessage:      sm,
	})

	return nil
}
