package darwintimetable

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/rest"
  "gopkg.in/robfig/cron.v2"
  "log"
  "time"
)

func (t *DarwinTimetable) ScheduleCleanup( c *cron.Cron ) {
  c.AddFunc( "0 0 2 * * *", func() {
    t.PruneSchedules()
  })
}

// PruneSchedules prunes all expired schedules
func (t *DarwinTimetable) PruneSchedules() ( int, error ) {
  count := 0

  if err := t.Update( func( tx *bolt.Tx ) error {
    lim := time.Now().Truncate( 24 * time.Hour )

    log.Println( "PruneSchedules:", lim.Format( "2006-01-02" ) )

    bucket := tx.Bucket( []byte( "DarwinJourney" ) )

    if err := bucket.ForEach( func( k, v []byte ) error {
      j := &Journey{}
      if j.fromBytes( v ) && j.SSD.Before( lim ) {
        count++
        return bucket.Delete( k )
      }
      return nil
    }); err != nil {
     return err
    }

    log.Println( "PruneSchedules:", count )
    return nil
  }); err != nil {
    return 0, err
  }
  return count, nil
}

func (dt *DarwinTimetable) PruneSchedulesHandler( r *rest.Rest ) error {
  if count, err := dt.PruneSchedules(); err != nil {
    return err
  } else {
    r.Status( 200 ).
    Value( count )
    return nil
  }
}
