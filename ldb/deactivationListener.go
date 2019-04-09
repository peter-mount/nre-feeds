package ldb

import (
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
)

// deactivationListener removes Services when a schedule is deactivated
func (d *LDB) deactivationListener(e *darwind3.DarwinEvent) {
	_ = d.Update(func(tx *bbolt.Tx) error {

		/*
			if e.Schedule != nil {
				// Ignore anything without a location & no Public times
				// Also ignore if the forecast says it's suppressed at this location.
				for _, l := range e.Schedule.Locations {
					if l.Times.IsPublic() {

						// Retrieve the station, it should be a valid one if we have Public times
						station := d.getStationTiploc(tx, l.Tiploc)
						if station != nil && station.removeService(tx, e.RID) {
							putStation(tx, station)
						}
					}
				}
			}
		*/

		if e.RID != "" {
			_ = removeSchedule(tx, e.RID)
		}

		return nil
	})
}
