package darwind3

// Process processes an inbound schedule importing or merging it with the
// current Schedule in the database
func (p *Schedule) Process( tx *Transaction ) error {

  old := tx.GetSchedule( p.RID )
  if old != nil {

    // If they are completely the same or the old entry is newer than the new one
    // then do nothing
    if p.Equals( old ) || tx.pport.TS.Before( old.Date ) {
      // They are identical so bail out
      return nil
    }

    // Use the new entry but merge in the locations from the old one so we keep
    // any forecasts
    ary := p.Locations

    // Run through old locations, any that match the new ones preserve the forecast
    for _, b := range old.Locations {
      for _, a := range ary {
        if a.EqualInSchedule( b ) {
          a.Forecast = b.Forecast
        }
      }
    }

    // Append any old locations not in the new one - we'll never remove a Location
    for _, b := range old.Locations {
      f := true
      for _, a := range ary {
        if a.EqualInSchedule( b ) {
          f = false
        }
      }
      if f {
        ary = append( ary, b)
      }
    }

    p.Locations = ary
  }

  return tx.PutSchedule( p )
}