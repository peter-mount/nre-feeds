package darwind3

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"strconv"
	"time"
)

// A location in a schedule.
// This is formed of the entries from a schedule and is updated by any incoming
// Forecasts.
//
// As schedules can be circular (i.e. start and end at the same station) then
// the unique key is Tiploc and CircularTimes.Time.
//
// Location's within a schedule are sorted by CircularTimes.Time accounting for
// crossing over midnight.
type Location struct {
	// Type of location, OR OPOR IP OPIP PP DT or OPDT
	Type string `json:"type"`
	// Tiploc of this location
	Tiploc string `json:"tiploc"`
	// The "display" time for this location
	// This is calculated using the first value in the following order:
	// Forecast.Time, Times.Time
	Time util.WorkingTime `json:"displaytime"`
	// The times for this entry
	Times util.CircularTimes `json:"timetable"`
	// TIPLOC of False Destination to be used at this location
	FalseDestination string `json:"falseDestination,omitempty"`
	// Is this service cancelled at this location
	Cancelled bool `json:"cancelled,omitempty"`
	// The Planned data for this location
	// i.e. information planned in advance
	Planned struct {
		// Current Activity Codes
		ActivityType string `json:"activity,omitempty"`
		// Planned Activity Codes (if different to current activities)
		PlannedActivity string `json:"plannedActivity,omitempty"`
		// A delay value that is implied by a change to the service's route.
		// This value has been added to the forecast lateness of the service at
		// the previous schedule location when calculating the expected lateness
		// of arrival at this location.
		RDelay int `json:"rDelay,omitempty"`
	} `json:"planned"`
	// The Forecast data at this location
	// i.e. information that changes in real time
	Forecast struct {
		// The "display" time for this location
		// This is calculated using the first value in the following order:
		// Departure, Arrival, Pass, or if none of those are set then the following
		// order in CircularTimes above is used: ptd, pta, wtd, wta & wtp
		Time util.WorkingTime `json:"time"`
		// If true then delayed. This is the delayed field in one of
		// Departure, Arrival, Pass in that order
		Delayed bool `json:"delayed,omitempty"`
		// If true then the train has arrived or passed this location
		Arrived bool `json:"arrived,omitempty"`
		// If true then the train has departed or passed this location
		Departed    bool `json:"departed,omitempty"`
		Passed      bool `json:"passed,omitempty"`
		Approaching bool `json:"approaching,omitempty"`
		// Forecast data for the arrival at this location
		Arrival util.TSTime `json:"arr,omitempty"`
		// Forecast data for the departure at this location
		Departure util.TSTime `json:"dep,omitempty"`
		// Forecast data for the pass of this location
		Pass util.TSTime `json:"pass,omitempty"`
		// Current platform number
		Platform Platform `json:"plat,omitempty"`
		// The service is suppressed at this location.
		Suppressed bool `json:"suppressed,omitempty"`
		// The length of the service at this location on departure
		// (or arrival at destination).
		// The default value of zero indicates that the length is unknown.
		Length int `json:"length,omitempty"`
		// Indicates from which end of the train stock will be detached.
		// The value is set to “true” if stock will be detached from the front of
		// the train at this location. It will be set at each location where stock
		// will be detached from the front.
		// Darwin will not validate that a stock detachment activity code applies
		// at this location.
		DetachFront bool `json:"detachFront,omitempty"`
		// The train order at this location (1, 2 or 3). 0 Means no TrainOrder has been set
		TrainOrder *TrainOrder `json:"trainOrder,omitempty"`
		// This is the TS time from Darwin when this Forecast was updated
		Date time.Time `json:"date,omitempty"`
	} `json:"forecast"`
	// The delay in seconds calculated as difference between forecast.time and timetable.time
	Delay int `json:"delay"`
	// Loading data for this location.
	Loading *Loading `json:"loading"`
	// updated if true trigger an event
	updated bool
}

// Compare compares two Locations by their times
func (a *Location) Compare(b *Location) bool {
	return b != nil && a.Times.Compare(&b.Times)
}

// Equals compares two Locations based on their Tiploc & working timetable.
// This is used when trying to locate a location that's been updated
func (a *Location) EqualInSchedule(b *Location) bool {
	return b != nil &&
		a.Tiploc == b.Tiploc &&
		a.Times.EqualInSchedule(&b.Times)
}

// Equals compares two Locations in their entirety
func (a *Location) Equals(b *Location) bool {
	return b != nil &&
		a.Type == b.Type &&
		a.Tiploc == b.Tiploc &&
		a.Times.Equals(&b.Times) &&
		a.FalseDestination == b.FalseDestination &&
		a.Cancelled == b.Cancelled &&
		a.Planned.ActivityType == b.Planned.ActivityType &&
		a.Planned.PlannedActivity == b.Planned.PlannedActivity &&
		a.Planned.RDelay == b.Planned.RDelay &&
		a.Forecast.Arrival.Equals(&b.Forecast.Arrival) &&
		a.Forecast.Departure.Equals(&b.Forecast.Departure) &&
		a.Forecast.Pass.Equals(&b.Forecast.Pass) &&
		a.Forecast.Platform.Equals(&b.Forecast.Platform) &&
		a.Forecast.Length == b.Forecast.Length &&
		a.Forecast.DetachFront == b.Forecast.DetachFront &&
		a.Forecast.TrainOrder == b.Forecast.TrainOrder
}

func (s *Location) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	// As we are unmarshalling from xml, mark the location as updated
	s.updated = true

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "tpl":
			s.Tiploc = attr.Value

		case "act":
			s.Planned.ActivityType = attr.Value

		case "planAct":
			s.Planned.PlannedActivity = attr.Value

		case "can":
			s.Cancelled = attr.Value == "true"

		case "fd":
			s.FalseDestination = attr.Value

		case "rdelay":
			if v, err := strconv.Atoi(attr.Value); err != nil {
				return err
			} else {
				s.Planned.RDelay = v
			}
		}
	}

	// Parse CircularTimes attributes
	s.Times.UnmarshalXMLAttributes(start)

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			var elem interface{}
			switch tok.Name.Local {
			case "arr":
				elem = &s.Forecast.Arrival

			case "dep":
				elem = &s.Forecast.Departure

			case "pass":
				elem = &s.Forecast.Pass

			case "plat":
				elem = &s.Forecast.Platform

			case "suppr":
				// TODO implement
				if err := decoder.Skip(); err != nil {
					return err
				}

			case "length":
				// TODO implement
				if err := decoder.Skip(); err != nil {
					return err
				}

			case "detachFront":
				// TODO implement
				if err := decoder.Skip(); err != nil {
					return err
				}

			default:
				if err := decoder.Skip(); err != nil {
					return err
				}
			}

			if elem != nil {
				if err := decoder.DecodeElement(elem, &tok); err != nil {
					return err
				}
			}

		case xml.EndElement:
			s.UpdateTime()
			return nil
		}
	}
}

// update sets all "calculated" fields
func (l *Location) UpdateTime() {
	l.Times.UpdateTime()

	if l.Loading != nil {
		l.Loading.Times.UpdateTime()
	}

	passed := l.Forecast.Pass.AT != nil && !l.Forecast.Pass.AT.IsZero()
	l.Forecast.Passed = passed

	l.Forecast.Departed = (l.Forecast.Departure.AT != nil && !l.Forecast.Departure.AT.IsZero()) || passed
	l.Forecast.Arrived = (l.Forecast.Arrival.AT != nil && !l.Forecast.Arrival.AT.IsZero()) || passed

	if l.Forecast.Departure.IsSet() {
		l.Forecast.Time = *l.Forecast.Departure.Time()
	} else if l.Forecast.Arrival.IsSet() {
		l.Forecast.Time = *l.Forecast.Arrival.Time()
	} else if l.Forecast.Pass.IsSet() {
		l.Forecast.Time = *l.Forecast.Pass.Time()
	} else if l.Times.Ptd != nil {
		l.Forecast.Time.Set(l.Times.Ptd.Get() * 60)
	} else if l.Times.Pta != nil {
		l.Forecast.Time.Set(l.Times.Pta.Get() * 60)
	} else if l.Times.Wtd != nil {
		l.Forecast.Time = *l.Times.Wtd
	} else if l.Times.Wta != nil {
		l.Forecast.Time = *l.Times.Wta
	} else if l.Times.Wtp != nil {
		l.Forecast.Time = *l.Times.Wtp
	} else {
		// Should never happen
		l.Forecast.Time.Set(-1)
	}

	// The Display time
	if l.Forecast.Time.IsZero() {
		l.Time = l.Times.Time
	} else {
		l.Time = l.Forecast.Time
	}

	l.Forecast.Delayed = l.Forecast.Departure.Delayed || l.Forecast.Arrival.Delayed || l.Forecast.Pass.Delayed

	if !l.Forecast.Time.IsZero() && !l.Times.Time.IsZero() {
		l.Delay = l.Forecast.Time.Get() - l.Times.Time.Get()
	} else {
		l.Delay = 0
	}

	var apt util.WorkingTime
	if l.Forecast.Arrival.ET != nil && l.Forecast.Arrival.ET.IsZero() {
		apt = *l.Forecast.Arrival.ET
	} else if l.Forecast.Pass.ET != nil && l.Forecast.Pass.ET.IsZero() {
		apt = *l.Forecast.Pass.ET
	}

	l.Forecast.Approaching = apt.IsApproaching()
}

// Clone makes a clone of a Location
func (a *Location) Clone() *Location {
	b := &Location{
		Type:             a.Type,
		Tiploc:           a.Tiploc,
		Times:            a.Times,
		FalseDestination: a.FalseDestination,
		Cancelled:        a.Cancelled,
		Planned:          a.Planned,
		Forecast:         a.Forecast,
		Loading:          a.Loading,
	}
	b.UpdateTime()
	return b
}

// MergeFrom merges data from one location into another
func (dest *Location) MergeFrom(src *Location) {
	dest.Times = src.Times

	// Copy the forecast but preserve metadata that won't be in the src
	trainOrder := dest.Forecast.TrainOrder
	dest.Forecast = src.Forecast
	dest.Forecast.TrainOrder = trainOrder

	// Mark location as updated
	dest.updated = true
}

func (l *Location) AddTiploc(m map[string]interface{}) {
	if l != nil {
		m[l.Tiploc] = nil
		if l.FalseDestination != "" {
			m[l.FalseDestination] = nil
		}
	}
}
