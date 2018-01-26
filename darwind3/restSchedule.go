package darwind3

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/rest"
)

func (d *DarwinD3) ScheduleHandler( r *rest.Rest ) error {
  return d.View( func( tx *bolt.Tx ) error {
    if sched := d.GetSchedule( tx, r.Var( "rid" ) ); sched != nil {
      sched.SetSelf( r )
      r.Status( 200 ).Value( sched )
    } else {
      r.Status( 404 )
    }

    return nil
  })
}
