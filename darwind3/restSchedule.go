package darwind3

import (
  "github.com/peter-mount/golib/rest"
)

func (d *DarwinD3) ScheduleHandler( r *rest.Rest ) error {
  if sched := d.GetSchedule( r.Var( "rid" ) ); sched != nil {
    sched.SetSelf( r )
    r.Status( 200 ).Value( sched )
  } else {
    r.Status( 404 )
  }

  return nil
}
