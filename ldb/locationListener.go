package ldb

import (
	"github.com/peter-mount/nre-feeds/darwind3"
)

// locationListener listens to location updates and updates the relevant
// Station with the new/updated entry
func (d *LDB) locationListener(e *darwind3.DarwinEvent) {
	e.Schedule.Update(func() error {

		// Ignore anything without a location & no Public times
		for idx, l := range e.Schedule.Locations {
			if l.Times.IsPublic() {

				// Retrieve the station, it should be a valid one if we have Public times
				station := d.GetStationTiploc(l.Tiploc)
				if station != nil {
					station.addService(e, idx)
				}
			}
		}

		return nil
	})
}
