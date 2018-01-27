package ldb

import (
  "darwind3"
  "log"
  "runtime"
)

// locationListener listens to location updates and updates the relevant
// Station with the new/updated entry
func (d *LDB) locationListener( c chan *darwind3.DarwinEvent ) {

  // As this listener has some heavy work to do we delegate to a pool of workers

  workers := 8 * runtime.NumCPU()
  log.Println( "LDB: Location listener", workers, "workers" )

  jobs := make( chan *darwind3.DarwinEvent, 1000 )
  for w := 1; w < workers; w++ {
    go d.locationListenerWorker( jobs )
  }

  for {
    e := <- c
    jobs <- e
  }
}

func (d *LDB) locationListenerWorker( c chan *darwind3.DarwinEvent ) {
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
