package darwintimetable

import (
  bolt "github.com/etcd-io/bbolt"
  "log"
  "time"
)

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
