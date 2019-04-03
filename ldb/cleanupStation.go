package ldb

import (
	"github.com/peter-mount/golib/statistics"
	"time"
)

// Cleans up a station removing old schedules
func (s *Station) cleanup() {
	now := time.Now()
	day := now.Add(-2 * time.Hour)

	s.Update(func() error {
		for rid, service := range s.services {
			if service.Timestamp().Before(day) {
				statistics.Incr("ldb.clean")
				delete(s.services, rid)
			}
		}
		return nil
	})
}

// Cleanup removes any old schedules still in memory for each station
func (st *Stations) Cleanup() {
	// Get set of current stations
	var crs []*Station
	st.Update(func() error {
		for _, s := range st.crs {
			crs = append(crs, s)
		}
		return nil
	})

	// Cleanup each one
	for _, s := range crs {
		s.cleanup()
	}
}
