package darwind3

// Process processes an inbound set of TrainOrders and applies them to the
// relevant schedules
func (to *trainOrderWrapper) Process( tx *Transaction ) error {

  // No order data then ignore
  if to.Set == nil {
    return nil
  }

  if to.Set.First != nil {
    if err := to.processOrder( tx, 1, to.Set.First ); err != nil {
      return err
    }
  }

  if to.Set.Second != nil {
    if err := to.processOrder( tx, 2, to.Set.Second ); err != nil {
      return err
    }
  }

  if to.Set.Third != nil {
    if err := to.processOrder( tx, 3, to.Set.Third ); err != nil {
      return err
    }
  }

  return nil
}

// Processes a specific TrainOrderItem
func (to *trainOrderWrapper) processOrder( tx *Transaction, order int, tod *trainOrderItem ) error {

  // Retrieve the schedule to be updated
  sched := tx.d3.GetSchedule( tod.RID )

  // No schedule then try to fetch it from the timetable
  if sched == nil {
    sched = tx.ResolveSchedule( tod.RID )
  }

  // Still no schedule then We've got a TS for a train with no known schedule so create one
  if sched == nil {
    sched = &Schedule{ RID: tod.RID }
    sched.Defaults()
  }

  // Locate the required location
  for _, l := range sched.Locations {
    if l.Tiploc == to.Tiploc && l.Times.Equals( &tod.Times ) {

      if to.Clear {
        l.Forecast.TrainOrder = nil
      } else {
        l.Forecast.TrainOrder = &TrainOrder{ Order: order, Platform: to.Platform }
      }

      // Mark as updated
      l.updated = true

      tx.d3.putSchedule( sched )
      return nil
    }
  }

  return nil
}
