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

func (d *LDB) OpenDB( dbFile string ) error {
  d.Stations = NewStations()

  // Add listeners
  d.Darwin.EventManager.ListenToEvents( darwind3.Event_LocationUpdated, d.locationListener )

  d.init()

  return nil
}

// init initialises the LDB database
func (d *LDB) init() {
  go d.initStations()
}
