package darwind3

import (
  bolt "github.com/coreos/bbolt"
)

type Transaction struct {
  // The root
  pport       *Pport
  //
  d3          *DarwinD3
  // Active Transaction
  tx          *bolt.Tx
  // Bucket for accessing the RID bucket
  ridBucket   *bolt.Bucket
}

// Processor interface used by some types used when processing a message and
// updating our internal state
type Processor interface {
  Process( *Transaction ) error
}

func (d *DarwinD3) ProcessUpdate( p *Pport, f func( *Transaction ) error ) error {
  return d.Update( func( tx *bolt.Tx ) error {
    t := &Transaction{
      pport: p,
      d3: d,
      tx: tx,
      ridBucket: tx.Bucket( []byte( "DarwinRID" ) ),
    }

    return f( t )
  })
}

// GetSchedule retrieves a schedule or nil if not found
func (tx *Transaction) GetSchedule( rid string ) *Schedule {
  sched := ScheduleFromBytes( tx.ridBucket.Get( []byte( rid ) ) )
  if sched == nil || sched.RID != rid {
    return nil
  }
  return sched
}

// PutSchedule persists a Schedule
func (tx *Transaction) PutSchedule( s *Schedule ) error {
  // Check existing entry & if the same dont persit
  existing := tx.GetSchedule( s.RID )
  if s.Equals( existing ) {
    return nil
  }

  if b, err := s.Bytes(); err != nil {
    return err
  } else {
    return tx.ridBucket.Put( []byte( s.RID ), b )
  }
}
