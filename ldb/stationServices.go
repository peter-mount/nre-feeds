package ldb

import (
	"github.com/peter-mount/nre-feeds/darwind3"
	"time"
)

const (
	dewllTime = time.Minute * 10
)

// Key must be unique so to support circular routes
// so we use the RID, tiploc and the timetable time
func (s *Station) key(sched *darwind3.Schedule, idx int) string {
	loc := sched.Locations[idx]
	return sched.RID + ":" + loc.Tiploc + ":" + loc.Times.Time.String()
}

// Adds a service to the station
func (s *Station) addService(e *darwind3.DarwinEvent, idx int) bool {
	// Only Public stations can be updated. Pass to the channel so the worker thread
	// can read it
	if s.Public && e.Schedule.Locations[idx].Times.IsPublic() {
		// Add service

		t := e.Schedule.GetTime(idx)
		if t.After(time.Now().Add(-dewllTime)) {
			k := s.key(e.Schedule, idx)

			service, exists := s.Services[k]
			if !exists {
				service = &Service{}
			}

			if service.update(e.Schedule, idx) {

				if !exists {
					s.Services[k] = service
					return true
				}

			}
		}

	}
	return false
}

func (s *Station) removeDepartedService(e *darwind3.DarwinEvent, idx int) bool {
	_, existed := s.Services[s.key(e.Schedule, idx)]
	return existed
}

// Removes all entries of a service from a station.
func (s *Station) removeService(rid string) {

	if s.Public {

		// As a service can call at a station more than once, scan all and remove
		// every instance of it.
		for k, service := range s.Services {
			if service.RID == rid {
				delete(s.Services, k)
			}
		}

	}

}
