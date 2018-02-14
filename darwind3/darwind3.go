// darwind3 handles the real time push port feed
package darwind3

type DarwinD3 struct {
  // Optional link to remote DarwinTimetable for resolving schedules.
  Timetable            string
  // Eventing
  EventManager         *DarwinEventManager
  // Schedule cache
  cache                 cache
  // Station message cache
  Messages             *StationMessages
}

// OpenDB opens a DarwinReference database.
func (r *DarwinD3) OpenDB( dbFile string ) error {
  r.EventManager = NewDarwinEventManager()
  r.Messages = NewStationMessages( dbFile )

  return r.cache.initCache( dbFile )
}
