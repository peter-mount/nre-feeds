package darwintimetable

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/rest"
)

func (dt *DarwinTimetable) JourneyHandler( r *rest.Rest ) error {
  return dt.View( func( tx *bolt.Tx ) error {
    if journey, exists := dt.GetJourney( tx, r.Var( "rid" ) ); exists {
      journey.SetSelf( r )
      r.Status( 200 ).Value( journey )
    } else {
      r.Status( 404 )
    }

    return nil
  })
}
