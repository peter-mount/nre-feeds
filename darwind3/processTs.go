package darwind3

// Process processes an inbound Train Status update, merging it with an existing
// schedule in the database
func (p *TS) Process( tx *Transaction ) error {

  // Retrieve the schedule to be updated
  sched := tx.GetSchedule( p.RID )

  // No schedule then try to fetch it from the timetable
  if sched == nil {
    sched = tx.ResolveSchedule( p.RID )
  }

  // Still no schedule then We've got a TS for a train with no known schedule so create one
  if sched == nil {
    sched = &Schedule{
      RID: p.RID,
      UID: p.UID,
      SSD: p.SSD,
    }
    sched.Defaults()
  }

  // If the TS is older than what's in the schedule then then do nothing as it's
  // presumably old data thats sent out of sync
  if tx.pport.TS.Before( sched.Date ) {
    return nil
  }

  // Run through schedule locations, any that match the new ones update the forecast
  for _, a := range sched.Locations {
    for _, b := range p.Locations {
      if a.EqualInSchedule( b ) {
        a.Times = b.Times
        a.Forecast = b.Forecast
      }
    }
  }

  // Append any locations not in the schedule
  for _, a := range p.Locations {
    f := true
    for _, b := range sched.Locations {
      if a.EqualInSchedule( b ) {
        f = false
      }
    }
    if f {
      sched.Locations = append( sched.Locations, a)
    }
  }

  // Finally sort the locations, set the date to that at Darwin then persist
  sched.Sort()
  sched.Date = tx.pport.TS

  return tx.PutSchedule( sched )
}
