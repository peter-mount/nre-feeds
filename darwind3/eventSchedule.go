package darwind3

// postScheduleEvent submits all events relating to a schedule
func (tx *Transaction) postScheduleUpdateEvent( s *Schedule ) {

  // Schedule updated
  tx.d3.EventManager.PostEvent( &DarwinEvent{
    Type: Event_ScheduleUpdated,
    RID: s.RID,
    Schedule: s,
  } )

}
