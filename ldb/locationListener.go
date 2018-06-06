package ldb

import (
  "github.com/peter-mount/nre-feeds/darwind3"
  "sort"
)

// locationListener listens to location updates and updates the relevant
// Station with the new/updated entry
func (d *LDB) locationListener( e *darwind3.DarwinEvent ) {
  e.Schedule.View( func() error {
    // Ensure the locations are in order
    loc := e.Schedule.Locations
    sort.SliceStable( loc, func( i, j int ) bool {
      return loc[ i ].Compare( loc[ j ] )
    } )

    // Ignore anything without a location & no public times
    for idx, l := range e.Schedule.Locations {
      if l.Times.IsPublic() {

        // Retrieve the station, it should be a valid one if we have public times
        station := d.GetStationTiploc( l.Tiploc )
        if station != nil {
          station.addService( e, idx )
        }
      }
    }

    return nil
  })
}
