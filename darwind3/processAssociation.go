package darwind3

import (
  "log"
)

// Process inbound associations
func (a *Association) Process( tx *Transaction ) error {
  log.Println( "Process Association", a )

  a.Date = tx.pport.TS

  err := a.Main.processSched( tx, a )
  if err != nil {
    return err
  }

  return a.Assoc.processSched( tx, a )
}

func (as *AssocService) processSched( tx *Transaction, a *Association ) error {
  sched := tx.d3.GetSchedule( as.RID )

  // No schedule then try to fetch it from the timetable
  if sched == nil {
    sched = tx.ResolveSchedule( as.RID )
  }

  // Still no schedule then We've got a TS for a train with no known schedule so create one
  if sched == nil {
    sched = &Schedule{ RID: as.RID }
    sched.Defaults()
  }

  if err := sched.Update( func() error {

    // Replace if we already have it
    for i, e := range sched.Associations {
      if e.Equals( a ) {
        sched.Associations[i] = a
        return nil
      }
    }

    // Not found then add it
    sched.Associations = append( sched.Associations, a )
    return nil
  }); err != nil {
    return err
  }

  sched.Date = tx.pport.TS
  tx.d3.putSchedule( sched )
  tx.d3.EventManager.PostEvent( &DarwinEvent{
    Type: Event_ScheduleUpdated,
    RID: sched.RID,
    Schedule: sched,
  })
  return nil
}
