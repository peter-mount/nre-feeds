// LDB - Live Departure Boards
package ldb

import (
	"github.com/peter-mount/nre-feeds/darwind3"
)

type LDB struct {
	// Link to D3
	Darwin string
	// Link to reference
	Reference string
	// Eventing
	EventManager *darwind3.DarwinEventManager
	// The managed stations
	Stations *Stations
}

func (d *LDB) Init() error {
	d.Stations = NewStations()

	// Add listeners
	d.EventManager.ListenToEvents(darwind3.Event_ScheduleUpdated, d.locationListener)
	d.EventManager.ListenToEvents(darwind3.Event_Deactivated, d.deactivationListener)
	d.EventManager.ListenToEvents(darwind3.Event_StationMessage, d.stationMessageListener)

	// init initialises the LDB memory structures to have the stations preloaded
	go d.initStations()

	return nil
}
