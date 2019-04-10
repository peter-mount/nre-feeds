package darwind3

import "github.com/etcd-io/bbolt"

// Process processes an inbound schedule importing or merging it with the
// current schedule in the database
func (p *Schedule) Process(tx *Transaction) error {
	if tx.d3.cache.tx != nil {
		return p.process(tx, tx.d3.cache.tx)
	}
	return tx.d3.Update(func(dbtx *bbolt.Tx) error {
		return p.process(tx, dbtx)
	})
}
func (p *Schedule) process(tx *Transaction, dbtx *bbolt.Tx) error {
	// Only look at an existing entry for uR messages. sR messages must replace the existing one
	if !tx.pport.SnapshotUpdate {
		old := GetSchedule(dbtx, p.RID)
		if old != nil {

			// If they are completely the same or the old entry is newer than the new one
			// then do nothing
			if p.Equals(old) || tx.pport.TS.Before(old.Date) {
				// They are identical so bail out
				return nil
			}

			// Use the new entry but merge in the locations from the old one so we keep
			// any forecasts
			ary := p.Locations

			// Run through old locations, any that match the new ones preserve the forecast
			for _, b := range old.Locations {
				for _, a := range ary {
					if a.EqualInSchedule(b) {
						a.Forecast = b.Forecast
					}
				}
			}

			// Append any old locations not in the new one - we'll never remove a Location
			for _, b := range old.Locations {
				f := true
				for _, a := range ary {
					if a.EqualInSchedule(b) {
						f = false
					}
				}
				if f {
					ary = append(ary, b)
				}
			}

			p.Locations = ary
		}
	}

	tx.d3.updateAssociations(dbtx, p)
	p.Date = tx.pport.TS
	p.Sort()
	if tx.d3.PutSchedule(p) {
		tx.d3.EventManager.PostEvent(&DarwinEvent{
			Type:     Event_ScheduleUpdated,
			RID:      p.RID,
			Schedule: p,
		})
	}
	return nil
}
