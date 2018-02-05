package ldb

import (
  "log"
  "time"
)

// Cleans up a station removing old schedules
func (s *Station) cleanup() int {
  ctr := 0
  now := time.Now()
  day := now.Add( -2 * time.Hour )

  s.Update( func() error {
    for rid, service := range s.services {
      if service.Timestamp().Before( day ) {
        ctr++
        delete( s.services, rid )
      }
    }
    return nil
  })

  return ctr
}

func (st *Stations) Cleanup() {
  // Get set of current stations
  var crs []*Station
  st.Update( func() error {
    for _, s := range st.crs {
      crs = append( crs, s )
    }
    return nil
  })

  // Cleanup each one
  ctr := 0
  for _, s := range crs {
    ctr += s.cleanup()
  }

  if ctr > 0 {
    log.Println( "LDB cleanup", ctr, "schedules removed")
  }
}
