package ldb

import (
  "darwind3"
)

// deactivationListener removes services when a schedule is deactivated
func (d *LDB) deactivationListener( e *darwind3.DarwinEvent ) {

  if e.Schedule != nil {
    e.Schedule.View( func() error {
      // Ignore anything without a location & no public times
      // Also ignore if the forecast says it's suppressed at this location.
      for _, l := range e.Schedule.Locations {
        if l.Times.IsPublic() {

          // Retrieve the station, it should be a valid one if we have public times
          station := d.GetStationTiploc( l.Tiploc )
          if station != nil {
            station.removeService( e.RID )
          }
        }
      }

      return nil
    })
  }

}
