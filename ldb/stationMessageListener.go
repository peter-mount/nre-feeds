package ldb

import (
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
)

func (d *LDB) RequestStationMessages() {
	d.EventManager.PostEvent(&darwind3.DarwinEvent{
		Type: darwind3.Event_Request_StationMessage,
	})
}

// deactivationListener removes Services when a schedule is deactivated
func (d *LDB) stationMessageListener(e *darwind3.DarwinEvent) {
	_ = d.Update(func(tx *bbolt.Tx) error {

		// Ensure all stations have this message
		if e.NewStationMessage != nil {
			if e.NewStationMessage.ID < 0 {
				// All stations get this message
				for _, s := range d.stations {
					//_ = tx.Bucket([]byte(crsBucket)).ForEach(func(k, v []byte) error {
					//s := StationFromBytes(v)
					if s != nil && s.Public {
						s.addStationMessage(e.NewStationMessage)
						d.putStation(tx, s)
					}
					//return nil
				} //)
			} else {
				// Only store in Public stations
				for _, crs := range e.NewStationMessage.Station {
					s := d.getStationCrs(tx, crs)
					if s != nil && s.Public {
						s.addStationMessage(e.NewStationMessage)
						d.putStation(tx, s)
					}
				}
			}
		}

		// Existing message so check for stations it no longer applies to
		if e.ExistingStationMessage != nil {
			m := make(map[string]int64)
			for _, crs := range e.ExistingStationMessage.Station {
				m[crs] = e.ExistingStationMessage.ID
			}

			// Remove crs codes that are present in the new message, this will leave
			// station's to delete
			if e.NewStationMessage != nil {
				for _, crs := range e.NewStationMessage.Station {
					delete(m, crs)
				}
			}

			for crs, id := range m {
				// Only store in Public stations
				s := d.getStationCrs(tx, crs)
				if s != nil && s.Public {
					updated := false

					var ary []int64

					for _, i := range s.Messages {
						if i != id {
							ary = append(ary, i)
						}
					}

					updated = len(s.Messages) != len(ary)
					s.Messages = ary

					if updated {
						d.putStation(tx, s)
					}
				}
			}
		}

		return nil
	})

}
