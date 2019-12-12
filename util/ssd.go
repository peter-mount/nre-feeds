package util

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

type SSD struct {
	t time.Time
}

func (a *SSD) Equals(b *SSD) bool {
	return b != nil && a.t == b.t
}

func (t *SSD) Parse(s string) {
	t.t, _ = time.Parse("2006-01-02", s)
}

// Before is an SSD before a specified time
func (s *SSD) Before(t time.Time) bool {
	return s.t.Before(t)
}

// Custom JSON Marshaler.
func (t *SSD) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *SSD) UnmarshalJSON(b []byte) error {
	s := string(b[:])
	if s != "null" && len(s) > 2 {
		t.Parse(s[1 : len(s)-1])
	}
	return nil
}

// Custom XML Marshaler.
func (t *SSD) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: t.String()}, nil
}

// String returns a SSD in "YYYY-MM-DD" format
func (t *SSD) String() string {
	return t.t.Format("2006-01-02")
}

func (t *SSD) Time() time.Time {
	return t.t.In(London())
}

func (t *SSD) Set(t0 time.Time) {
	t.t = t0
}
