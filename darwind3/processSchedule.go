package darwind3

import (
  "log"
)

// Processor interface
func (p *Schedule) Process( tx *Transaction ) error {
  log.Println( p )

  old := tx.GetSchedule( p.RID )
  if old == nil {
    // New schedule so simply persist
    return tx.PutSchedule( p )
  }

  log.Println( "Found existing", old )

  // Use the new entry but merge in the locations from the old one so we keep
  // any forecasts
  

  if p.Equals( old ) {
    // They are identical so bail out
    log.Println( "Bailing as equal")
    return nil
  }

  // Merge into old & persist
  return tx.PutSchedule( old )
}
