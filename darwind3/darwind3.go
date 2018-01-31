// darwind3 handles the real time push port feed
package darwind3

import (
  "darwintimetable"
)

type DarwinD3 struct {
  // Optional link to DarwinTimetable for resolving schedules.
  Timetable            *darwintimetable.DarwinTimetable
  // Eventing
  EventManager         *DarwinEventManager
  // Schedule cache
  cache                 cache
}

// OpenDB opens a DarwinReference database.
func (r *DarwinD3) OpenDB( dbFile string ) error {
  r.EventManager = NewDarwinEventManager()

  return r.cache.initCache( dbFile )
}
