package util

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
)

// Public Timetable time
// Note: 00:00 is not possible as in CIF that means no-time
type PublicTime struct {
	t int
	n bool
}

func (a *PublicTime) Equals(b *PublicTime) bool {
	if a == nil {
		return b == nil || b.IsZero()
	}
	if b == nil {
		return a.IsZero()
	}
	if a.IsZero() {
		return b.IsZero()
	}
	return a.t == b.t
}

// Compare a PublicTime against another, accounting for crossing midnight.
// The rules for handling crossing midnight are:
// < -6 hours = crossed midnight
// < 0 back in time
// < 18 hours increasing time
// > 18 hours back in time & crossing midnight
func (a *PublicTime) Compare(b *PublicTime) bool {
	if b == nil {
		return false
	}

	d := a.t - b.t

	if d < -360 || d > 1080 {
		return d > 0
	}

	return d < 0
}

// NewPublicTime returns a new PublicTime instance from a string of format "HH:MM"
func NewPublicTime(s string) *PublicTime {
	v := &PublicTime{}
	v.Parse(s)
	return v
}

func (v *PublicTime) Parse(s string) {
	if s == "" {
		v.t = -1
	} else {
		a, _ := strconv.Atoi(s[0:2])
		b, _ := strconv.Atoi(s[3:5])
		v.Set((a * 60) + b)
	}
}

// Custom JSON Marshaler. This will write null or the time as "HH:MM"
func (t *PublicTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(t.String())
}

func (t *PublicTime) UnmarshalJSON(b []byte) error {
	s := string(b[:])
	if s != "null" && len(s) > 2 {
		t.Parse(s[1 : len(s)-1])
	}
	return nil
}

// Custom XML Marshaler.
func (t *PublicTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if t.IsZero() {
		return xml.Attr{}, nil
	}
	return xml.Attr{Name: name, Value: t.String()}, nil
}

// String returns a PublicTime in HH:MM format or 5 blank spaces if it's not set.
func (t *PublicTime) String() string {
	if t.IsZero() {
		return "     "
	}

	return fmt.Sprintf("%02d:%02d", t.t/60, t.t%60)
}

// Get returns the PublicTime in minutes of the day
func (t *PublicTime) Get() int {
	return t.t
}

// Set sets the PublicTime in minutes of the day
func (t *PublicTime) Set(v int) {
	t.t = v
}

// IsZero returns true if the time is not present
func (t *PublicTime) IsZero() bool {
	return t.t <= 0
}

// Is this instance representing nil
func (t *PublicTime) IsNil() bool {
	return t.n
}
