package util

import (
	"encoding/xml"
	"fmt"
)

// A scheduled time used to distinguish a location on circular routes.
// Note that all scheduled time attributes are marked as optional,
// but at least one must always be supplied.
// Only one value is required, and typically this should be the wtd value.
// However, for locations that have no wtd, or for clients that deal
// exclusively with public times, another value that is valid for the
// location may be supplied.
type CircularTimes struct {
	// The time for this location.
	// This is calculated as the first value defined below in the following
	// sequence: Wtd, Wta, Wtp, Ptd & Pta.
	Time WorkingTime `json:"time"`
	// Public Scheduled Time of Arrival
	Pta *PublicTime `json:"pta,omitempty"`
	// Public Scheduled Time of Departure
	Ptd *PublicTime `json:"ptd,omitempty"`
	// Working Scheduled Time of Arrival
	Wta *WorkingTime `json:"wta,omitempty"`
	// Working Scheduled Time of Departure
	Wtd *WorkingTime `json:"wtd,omitempty"`
	// Working Scheduled Time of Passing
	Wtp *WorkingTime `json:"wtp,omitempty"`
}

// IsPublic returns true of the instance contains public times
func (t *CircularTimes) IsPublic() bool {
	return t.Pta != nil || t.Ptd != nil
}

// IsPass returns true if the instance represents a pass at a station
func (t *CircularTimes) IsPass() bool {
	return t.Wtp != nil
}

// UnmarshalXMLAttributes reads from an arbitary start element
func (t *CircularTimes) UnmarshalXMLAttributes(start xml.StartElement) {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "pta":
			t.Pta = NewPublicTime(attr.Value)

		case "ptd":
			t.Ptd = NewPublicTime(attr.Value)

		case "wta":
			t.Wta = NewWorkingTime(attr.Value)

		case "wtd":
			t.Wtd = NewWorkingTime(attr.Value)

		case "wtp":
			t.Wtp = NewWorkingTime(attr.Value)
		}
	}
	t.UpdateTime()
}

// Compare compares two Locations by their times
func (a *CircularTimes) Compare(b *CircularTimes) bool {
	return b != nil && a.Time.Compare(&b.Time)
}

// UpdateTime updates the Time field used for sequencing the location.
// This is the the first one of these set in the following order:
// Wtd, Wta, Wtp, Ptd, Pta
// Note this value is not persisted as it's a generated value
func (l *CircularTimes) UpdateTime() {
	t := -1

	if l.Wtd != nil && !l.Wtd.IsZero() {
		t = l.Wtd.Get()
	} else if l.Wta != nil && !l.Wta.IsZero() {
		t = l.Wta.Get()
	} else if l.Wtp != nil && !l.Wtp.IsZero() {
		t = l.Wtp.Get()
	} else if l.Ptd != nil && !l.Ptd.IsZero() {
		// Should not happen, we should have a working time
		t = l.Ptd.Get() * 60
	} else if l.Pta != nil && !l.Pta.IsZero() {
		// Should not happen, we should have a working time
		t = l.Pta.Get() * 60
	}

	l.Time.Set(t)
}

// Equals returns true if both CircularTimes are exactly the same
func (a *CircularTimes) Equals(b *CircularTimes) bool {
	return b != nil &&
		a.Pta.Equals(b.Pta) &&
		a.Ptd.Equals(b.Ptd) &&
		a.Wta.Equals(b.Wta) &&
		a.Wtd.Equals(b.Wtd) &&
		a.Wtp.Equals(b.Wtp)
}

// EqualInSchedule returns true if the working timetable fields of both
// CircularTimes are equal as they are the primary key
func (a *CircularTimes) EqualInSchedule(b *CircularTimes) bool {
	if a == nil {
		return b == nil
	}
	return b != nil &&
		a.Wta.Equals(b.Wta) &&
		a.Wtd.Equals(b.Wtd) &&
		a.Wtp.Equals(b.Wtp)
}

func (l *CircularTimes) String() string {
	return fmt.Sprintf("%8v %5v %5v %8v %8v %8v", &l.Time, l.Pta, l.Ptd, l.Wta, l.Wtd, l.Wtp)
}
