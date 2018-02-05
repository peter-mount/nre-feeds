package ldb

import (
  "darwind3"
)

// deactivationListener removes services when a schedule is deactivated
func (d *LDB) stationMessageListener( c chan *darwind3.DarwinEvent ) {
  for {
    e := <- c

    // Ensure all stations have this message
    if e.NewStationMessage != nil {
      for _, crs := range e.NewStationMessage.Station {
        // Only store in public stations
        if s := d.GetStationCrs( crs ); s != nil && s.public {
          updated := false
          s.Update( func() error {
            for _, i := range s.messages {
              if i == e.NewStationMessage.ID {
                return nil
              }
            }
            s.messages = append( s.messages, e.NewStationMessage.ID )
            updated = true
            return nil
          })

          if updated {
            s.update()
          }
        }
      }
    }

    // Existing message so check for stations it no longer applies to
    if e.ExistingStationMessage != nil {
      var m map[string]int = make( map[string]int )
      for _, crs := range e.ExistingStationMessage.Station {
        m[ crs ] = e.ExistingStationMessage.ID
      }

      // Remove crs codes that are present in the new message, this will leave
      // station's to delete
      if e.NewStationMessage != nil {
        for _, crs := range e.NewStationMessage.Station {
          delete( m, crs )
        }
      }

      for crs, id := range m {
        // Only store in public stations
        if s := d.GetStationCrs( crs ); s != nil && s.public {
          updated := false
          s.Update( func() error {
            var ary []int

            for _, i := range s.messages {
              if i != id {
                ary = append( ary, i )
              }
            }

            updated = len( s.messages ) != len( ary )
            s.messages = ary
            return nil
          })

          if updated {
            s.update()
          }
        }
      }
    }
  }
}
