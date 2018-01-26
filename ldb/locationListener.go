package ldb

import (
  "darwind3"
)

// locationListener listens to location updates and updates the relevant
// Station with the new/updated entry
/*
func (d *Departures) locationListener1( c chan *darwind3.DarwinEvent ) {
  // fire off 4 workers on this channel
  go d.locationListenerWorker( c )
  go d.locationListenerWorker( c )
  go d.locationListenerWorker( c )
  // Last one runs in this thread
  d.locationListenerWorker( c )
}*/

func (d *LDB) locationListener( c chan *darwind3.DarwinEvent ) {
  for {
    e := <- c

    // Ignore anything without a location & no public times
    // Also ignore if the forecast says it's suppressed at this location.
    if e.Location != nil &&
       e.Location.Times.IsPublic() &&
       !e.Location.Forecast.Suppressed {

      // Retrieve the station, it should be a valid one if we have public times
      station := d.GetStationTiploc( e.Location.Tiploc )
      if station != nil && station.public {
        station.addService( e )
      }
    }
  }
}
