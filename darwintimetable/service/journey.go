package service

import (
  bolt "github.com/etcd-io/bbolt"
  "github.com/peter-mount/golib/rest"
)

func (dt *DarwinTimetableService) JourneyHandler( r *rest.Rest ) error {
  return dt.timetable.View( func( tx *bolt.Tx ) error {
    if journey, exists := dt.timetable.GetJourney( tx, r.Var( "rid" ) ); exists {
      journey.SetSelf( r )
      r.Status( 200 ).
        JSON().
        Value( journey )
    } else {
      r.Status( 404 )
    }

    return nil
  })
}
