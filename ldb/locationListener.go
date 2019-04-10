package ldb

import (
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
)

// locationListener listens to location updates and updates the relevant
// Station with the new/updated entry
func (d *LDB) locationListener(e *darwind3.DarwinEvent) {
	_ = d.Update(func(tx *bbolt.Tx) error {

		// Ignore anything without a location & no Public times
		for idx, l := range e.Schedule.Locations {
			if l.Times.IsPublic() {

				// Retrieve the station, it should be a valid one if we have Public times
				station := d.getStationTiploc(tx, l.Tiploc)
				if station != nil && station.Public {
					//updated := false

					if l.Forecast.Departed {
						_ = station.removeDepartedService(tx, e, idx)
					} else {
						_ = station.addService(tx, e, idx)
					}

					/*
						if updated {
							putStation(tx, station)
						}
					*/

				}
			}
		}

		_ = darwind3.PutSchedule(tx, e.Schedule)

		return nil
	})
}
