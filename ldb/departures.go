// LDB - Live Departure Boards
package ldb

import (
  "darwind3"
  "darwinref"
)

type LDB struct {
  // Link to D3
  Darwin       *darwind3.DarwinD3
  // Link to reference
  Reference    *darwinref.DarwinReference
  // The managed stations
  Stations     *Stations
}

func (d *LDB) Init() error {
  d.Stations = NewStations()

  // Add listeners
  d.Darwin.EventManager.ListenToEventsCapacity( darwind3.Event_ScheduleUpdated, 10000, d.locationListener )
  d.Darwin.EventManager.ListenToEventsCapacity( darwind3.Event_Deactivated, 10000, d.deactivationListener )

  // init initialises the LDB memory structures to have the stations preloaded
  go d.initStations()

  return nil
}
