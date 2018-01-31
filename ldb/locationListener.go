package ldb

import (
  "darwind3"
)

// locationListener listens to location updates and updates the relevant
// Station with the new/updated entry
func (d *LDB) locationListener( c chan *darwind3.DarwinEvent ) {
  for {
    e := <- c

    e.Schedule.View( func() error {
      // Ignore anything without a location & no public times
      // Also ignore if the forecast says it's suppressed at this location.
      for _, l := range e.Schedule.Locations {
        if l.Times.IsPublic() && !l.Forecast.Suppressed {

          // Retrieve the station, it should be a valid one if we have public times
          station := d.GetStationTiploc( l.Tiploc )
          if station != nil {
            station.addService( e, l )
          }
        }
      }

      return nil
    })
  }
}
