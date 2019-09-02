package darwind3

import (
  "github.com/peter-mount/nre-feeds/util"
  "strings"
)

// A calling point of a service after a station.
type CallingPoint struct {
  // Tiploc of this location
  Tiploc      string           `json:"tpl"`
  Time        util.WorkingTime `json:"time"`
  Delay       int              `json:"delay"`
  Delayed     bool             `json:"delayed,omitempty"`
  Approaching bool             `json:"approaching,omitempty"`
  At          bool             `json:"at,omitempty"`
  Departed    bool             `json:"departed,omitempty"`
  Passed      bool             `json:"passed,omitempty"`
  SetDownOnly bool             `json:"setDownOnly,omitempty"`
}

// IsCallingPoint returns true if this location is a valid CallingPoint.
// For a Location to be a CallingPoint, it has to be not-cancelled,
// not passed (arrived||departed) and have a public departure time
// (or arrival for destination only)
func (l *Location) IsCallingPoint() bool {
  // Not cancelled, arrived nor departed
  if l.Cancelled || l.Forecast.Arrived || l.Forecast.Departed || l.Forecast.Suppressed {
    return false
  }

  // For destination only, has a valid Public Arrival Time
  // nre-feeds#17 - Allow set-down only stops to be included
  if (l.IsDestination() || l.IsSetDownOnly()) && l.HasPublicArrival() {
    return true
  }

  // Only valid for Departure times
  return l.HasPublicDeparture()
}

func (l *Location) AsCallingPoint() CallingPoint {
  l.UpdateTime()
  return CallingPoint{
    Tiploc:      l.Tiploc,
    Time:        l.Time,
    Delay:       l.Delay,
    Delayed:     l.Forecast.Delayed,
    At:          !l.Forecast.Passed && l.Forecast.Arrived && !l.Forecast.Departed,
    Departed:    !l.Forecast.Passed && l.Forecast.Departed,
    Passed:      l.Forecast.Passed,
    Approaching: l.Forecast.Time.IsApproaching(),
    SetDownOnly: l.IsSetDownOnly(),
  }
}

// HasPublicArrival returns true if this Location has a public timetable arrival time.
func (l *Location) HasPublicArrival() bool {
  return l != nil && l.Times.Pta != nil && !l.Times.Pta.IsZero()
}

// HasPublicDeparture returns true if this Location has a public timetable departure time.
func (l *Location) HasPublicDeparture() bool {
  return l != nil && l.Times.Ptd != nil && !l.Times.Ptd.IsZero()
}

// IsOrigin returns true if this Location is the train's origin.
// Origin's are Location's with types OR or OPOR.
func (l *Location) IsOrigin() bool {
  return l != nil && (l.Type == "OR" || l.Type == "OPOR")
}

// IsDestination returns true if this Location is the train's destination.
// Destination's are Location's with types DT or OPDT.
func (l *Location) IsDestination() bool {
  return l != nil && (l.Type == "DT" || l.Type == "OPDT")
}

// IsSetDownOnly returns true if this Location is an IP that allows passengers to set-down only & not board the train.
// SetDown only Locations are those of type IP and contain Activity D.
func (l *Location) IsSetDownOnly() bool {
  return l != nil && l.Type == "IP" && ContainsActivity(l.Planned.ActivityType, "D ")
}

// ContainsActivity returns true if the activity list contains the required activity code.
// Note, required must be a 2 character string. If a single character code is provided then it must have a trailing
// space. e.g. For "D" (Set down only) then required must be set to "D " - note the trailing space.
func ContainsActivity(activity, required string) bool {
  return activity != "" && len(required) == 2 && (strings.Index(activity, required)%2) == 0
}

func (l *Location) After(other *Location) bool {
  if l == nil || other == nil {
    return false
  }

  diff := l.Time.Difference(other.Time)

  return diff > 0
}

// GetCallingPoints returns a list of calling points from a specific location
// in the schedule. If the specific location has a FalseDestination set then
// the calling point list will terminate there rather than at the
// end of the schedule.
func (s *Schedule) GetCallingPoints(idx int) []CallingPoint {
  var cp []CallingPoint

  if idx >= 0 && (idx+1) < len(s.Locations) {
    loc := s.Locations[idx]

    // departureboards#4 If set then we need to stop at this tiploc rather than the entire schedule
    falseDest := loc.FalseDestination

    for _, l := range s.Locations[idx+1:] {
      // Filter by calling point & if it's after the location time else during delays
      // we see the previous entries first
      if l.IsCallingPoint() && l.Time.After(&loc.Time) && l.Tiploc != loc.Tiploc {
        cp = append(cp, l.AsCallingPoint())
      }

      // departureboards#4 Stop at the falseDest (if one is defined)
      if l.Tiploc == falseDest {
        return cp
      }
    }
  }

  return cp
}

// GetLastReport returns the last report as a CallingPoint
func (s *Schedule) GetLastReport() CallingPoint {

  var cp *Location
  for _, l := range s.Locations {
    l.UpdateTime()
    if (l.Forecast.Arrived || l.Forecast.Departed || l.Forecast.Time.IsApproaching()) && !l.Cancelled {
      cp = l
    }
  }
  if cp != nil {
    return cp.AsCallingPoint()
  }
  return CallingPoint{}
}
