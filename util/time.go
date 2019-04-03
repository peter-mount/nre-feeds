package util

import (
	"time"
)

const (
	DateTime      = "2006-01-02 15:04:05"
	Date          = "2006-01-02"
	HumanDateTime = "2006 Jan 02 15:04:05"
	HumanDate     = "2006 Jan 02"
	Time          = "15:04:05"
)

func London() *time.Location {
	l, _ := time.LoadLocation("Europe/London")
	return l
}

// Now is the same as time.Now() but the returned value will be in the
// Europe/London timezone.
func Now() time.Time {
	return time.Now().In(London())
}

// Time returns the PublicTime as a time.Time
// The returned value will be in the Europe/London time zone.
// t The time containing the date to contain the PublicTime
func (pt *PublicTime) Time(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour).In(London()).Add(time.Duration(pt.Get()) * time.Minute)
}

// TrainTime returns the PublicTime as a time.Time based on the supplied
// start time. This start time should be the date & time of the scheduled
// departure because we use this moment in time to determine if midnight has
// been crossed so that the result is always after the passed value of t.
// The returned value will be in the Europe/London time zone.
// t The time containing the date & time of the scheduled departure.
func (pt *PublicTime) TrainTime(t time.Time) time.Time {
	var r = pt.Time(t)
	if r.Before(t) {
		r = r.Add(24 * time.Hour)
	}
	return r
}

// Time returns the WorkingTime as a time.Time.
// The returned value will be in the Europe/London time zone.
// t The time containing the date to contain the WorkingTime
func (wt *WorkingTime) Time(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour).In(London()).Add(time.Duration(wt.Get()) * time.Second)
}

// TrainTime returns the PublicTime as a time.Time based on the supplied
// start time. This start time should be the date & time of the scheduled
// departure because we use this moment in time to determine if midnight has
// been crossed so that the result is always after the passed value of t.
// The returned value will be in the Europe/London time zone.
// t The time containing the date & time of the scheduled departure.
func (wt *WorkingTime) TrainTime(t time.Time) time.Time {
	var r = wt.Time(t)
	if r.Before(t) {
		r = r.Add(24 * time.Hour)
	}
	return r
}
