package darwind3

import (
  "github.com/peter-mount/nre-feeds/darwintimetable"
  "github.com/peter-mount/nre-feeds/darwintimetable/client"
  "github.com/peter-mount/nre-feeds/util"
  "strconv"
)

// ResolveSchedule attempts to retrieve a schedule from the timetable.
// If DarwinD3.Timetable is not set then this always returns nil
func (d *Transaction) ResolveSchedule( rid string ) *Schedule {
  if d.d3.Timetable == "" {
    return nil
  }

  client := &client.DarwinTimetableClient{ Url: d.d3.Timetable }

  if journey, err := client.GetJourney( rid ); err != nil {
    return nil
  } else if journey == nil {
    return nil
  } else {

    s := &Schedule{
      RID: journey.RID,
      UID: journey.UID,
      TrainId: journey.TrainID,
      SSD: journey.SSD,
      Toc: journey.Toc,
      TrainCat: journey.TrainCat,
      PassengerService: journey.Passenger,
      Active: true,
      // TODO check this is correct if j.Passenger is false
      Status: "P",
    }

    s.CancelReason.Reason = journey.CancelReason

    // Import locations
    for _, tl := range journey.Schedule {
      l := &Location{
        Type: tl.Type,
        Tiploc: tl.Tiploc,
        FalseDestination: tl.FalseDest,
        Cancelled: tl.Cancelled,
      }

      l.Times.Pta = tl.Pta
      l.Times.Ptd = tl.Ptd
      l.Times.Wta = tl.Wta
      l.Times.Wtd = tl.Wtd
      l.Times.Wtp = tl.Wtp

      l.Planned.ActivityType = tl.Act
      l.Planned.PlannedActivity = tl.PlanAct

      if tl.RDelay != "" {
        if d, e := strconv.Atoi( tl.RDelay ); e == nil {
          l.Planned.RDelay = d
        }
      }

      l.Forecast.Platform.Platform = tl.Platform

      s.Locations = append( s.Locations, l )
    }

    // Import associations
    for _, ta := range journey.Associations {
      a := &Association{
        Tiploc: ta.Tiploc,
        Category: ta.Category,
        Cancelled: ta.Cancelled,
        Deleted: ta.Deleted,
        Date: ta.Date,
      }

      a.Main.copyFromTimetable( &ta.Main )
      a.Assoc.copyFromTimetable( &ta.Assoc )

      s.Associations = append( s.Associations, a )
    }

    s.Sort()

    return s
  }
}

func (a *AssocService) copyFromTimetable( b *darwintimetable.AssocService ) {
  a.RID = b.RID

  if b.Wta != "" {
    a.Times.Wta = util.NewWorkingTime( b.Wta )
  }

  if b.Wtd != "" {
    a.Times.Wtd = util.NewWorkingTime( b.Wtd )
  }

  if b.Wtp != "" {
    a.Times.Wtp = util.NewWorkingTime( b.Wtp )
  }

  if b.Pta != "" {
    a.Times.Pta = util.NewPublicTime( b.Pta )
  }

  if b.Ptd != "" {
    a.Times.Ptd = util.NewPublicTime( b.Ptd )
  }
}
