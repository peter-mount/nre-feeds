package darwind3

import (
  "util"
)

// A calling point of a service after a station.
type CallingPoint struct {
  Tiploc  string            `json:"tpl"`
  Time    util.WorkingTime  `json:"time"`
}

// IsCallingPoint returns true if this location is a valid CallingPoint.
// For a Location to be a CallingPoint, it has to be not-cancelled,
// not passed (arrived||departed) and have a public departure time
// (or arrival for destination only)
func (l *Location) IsCallingPoint() bool {
  // Not cancelled, arrived nor departed
  if l.Cancelled || l.Forecast.Arrived || l.Forecast.Departed {
    return false
  }

  // For destination only, has a valid Public Arrival Time
  if ( l.Type == "DT" || l.Type == "OPDT" ) &&
      l.Times.Pta != nil && !l.Times.Pta.IsZero() {
    return true
  }

  // Only valid for Departure times
  return l.Times.Ptd != nil && !l.Times.Ptd.IsZero()
}

// GetCallingPoints returns a list of calling points from a specific location
// in the schedule.
func (s *Schedule) GetCallingPoints( idx int ) []*CallingPoint {
  var cp []*CallingPoint

  if idx >= 0 && (idx+1) < len( s.Locations ) {
    for _, l := range s.Locations[ idx+1: ] {
      if l.IsCallingPoint() {
        cp = append( cp, &CallingPoint{ Tiploc: l.Tiploc, Time: l.Time } )
      }
    }
  }

  return cp
}
