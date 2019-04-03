package ldb

import (
	"github.com/peter-mount/nre-feeds/darwind3"
)

// Adds a service to the station
func (s *Station) addService(e *darwind3.DarwinEvent, idx int) {
	// Only public stations can be updated. Pass to the channel so the worker thread
	// can read it
	if s.public && e.Schedule.Locations[idx].Times.IsPublic() {

		s.Update(func() error {

			loc := e.Schedule.Locations[idx]

			k := e.RID + ":" + loc.Times.Time.String()

			service, exists := s.services[k]
			if !exists {
				service = &Service{}
			}

			if service.update(e.Schedule, idx) {

				if !exists {
					// Key must be unique so to support circular routes
					// so we use the RID, tiploc and the timetable time
					k := e.RID + ":" + loc.Tiploc + ":" + loc.Times.Time.String()
					s.services[k] = service
				}

				s.update()

			}

			return nil
		})

	}
}

// Removes a service from a station.
func (s *Station) removeService(rid string) {

	if s.public {

		s.Update(func() error {

			// As a service can call at a station more than once, scan all and remove
			// every instance of it.
			for k, service := range s.services {
				if service.RID == rid {
					delete(s.services, k)
				}
			}

			return nil
		})
	}

}
