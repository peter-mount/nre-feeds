package darwind3

import (
	"github.com/peter-mount/nre-feeds/util"
)

// A calling point of a service after a station.
type CallingPoint struct {
	// Tiploc of this location
	Tiploc string `json:"tpl"`
	// Time due at this location
	Time util.WorkingTime `json:"time"`
	// Delay in seconds
	Delay int `json:"delay"`
	// Delayed at this location
	Delayed bool `json:"delayed"`
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
	if (l.Type == "DT" || l.Type == "OPDT") &&
		l.Times.Pta != nil && !l.Times.Pta.IsZero() {
		return true
	}

	// Only valid for Departure times
	return l.Times.Ptd != nil && !l.Times.Ptd.IsZero()
}

func (l *Location) AsCallingPoint() *CallingPoint {
	return &CallingPoint{
		Tiploc:  l.Tiploc,
		Time:    l.Time,
		Delay:   l.Delay,
		Delayed: l.Forecast.Delayed,
	}
}

// GetCallingPoints returns a list of calling points from a specific location
// in the schedule. If the specific location has a FalseDestination set then
// the calling point list will terminate there rather than at the
// end of the schedule.
func (s *Schedule) GetCallingPoints(idx int) []*CallingPoint {
	var cp []*CallingPoint

	if idx >= 0 && (idx+1) < len(s.Locations) {
		// departureboards#4 If set then we need to stop at this tiploc rather than the entire schedule
		falseDest := s.Locations[idx].FalseDestination

		for _, l := range s.Locations[idx+1:] {
			if l.IsCallingPoint() {
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
func (s *Schedule) GetLastReport() *CallingPoint {
	var cp *Location
	for _, l := range s.Locations {
		if l.Forecast.Arrived || l.Forecast.Departed {
			cp = l
		}
	}
	if cp != nil {
		return cp.AsCallingPoint()
	}
	return nil
}
