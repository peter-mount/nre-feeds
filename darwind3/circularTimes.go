package darwind3

import (
  "darwintimetable"
  "fmt"
  "github.com/peter-mount/golib/codec"
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
  Time              darwintimetable.WorkingTime
  // Public Scheduled Time of Arrival
  Pta              *darwintimetable.PublicTime
  // Public Scheduled Time of Departure
  Ptd              *darwintimetable.PublicTime
  // Working Scheduled Time of Arrival
  Wta              *darwintimetable.WorkingTime
  // Working Scheduled Time of Departure
  Wtd              *darwintimetable.WorkingTime
  // Working Scheduled Time of Passing
  Wtp              *darwintimetable.WorkingTime
}

// Compare compares two Locations by their times
func (a *CircularTimes) Compare( b *CircularTimes ) bool {
  return b != nil && a.Time.Compare( &b.Time )
}

// UpdateTime updates the Time field used for sequencing the location.
// This is the the first one of these set in the following order:
// Wtd, Wta, Wtp, Ptd, Pta
// Note this value is not persisted as it's a generated value
func (l *CircularTimes) UpdateTime() {
  t := -1

  if l.Wtd != nil {
    t = l.Wtd.Get()
  } else if l.Wta != nil {
    t = l.Wta.Get()
  } else if l.Wtp != nil {
    t = l.Wtp.Get()
  } else if l.Ptd != nil {
    // Should not happen, we should have a working time
    t = l.Ptd.Get() * 60
  } else if l.Ptd != nil {
    // Should not happen, we should have a working time
    t = l.Pta.Get() * 60
  }

  l.Time.Set( t )
}

func (a *CircularTimes) Equals( b *CircularTimes ) bool {
  return b != nil &&
    darwintimetable.PublicTimeEquals( a.Pta, b.Pta ) &&
    darwintimetable.PublicTimeEquals( a.Ptd, b.Ptd ) &&
    darwintimetable.WorkingTimeEquals( a.Wta, b.Wta ) &&
    darwintimetable.WorkingTimeEquals( a.Wtd, b.Wtd ) &&
    darwintimetable.WorkingTimeEquals( a.Wtp, b.Wtp )
}

func (t *CircularTimes) Write( c *codec.BinaryCodec ) {
  darwintimetable.PublicTimeWrite( c, t.Pta )
  darwintimetable.PublicTimeWrite( c, t.Ptd )
  darwintimetable.WorkingTimeWrite( c, t.Wta )
  darwintimetable.WorkingTimeWrite( c, t.Wtd )
  darwintimetable.WorkingTimeWrite( c, t.Wtp )
}

func (t *CircularTimes) Read( c *codec.BinaryCodec ) {
  t.Pta = darwintimetable.PublicTimeRead( c )
  t.Ptd = darwintimetable.PublicTimeRead( c )
  t.Wta = darwintimetable.WorkingTimeRead( c )
  t.Wtd = darwintimetable.WorkingTimeRead( c )
  t.Wtp = darwintimetable.WorkingTimeRead( c )
  t.UpdateTime()
}

func (l *CircularTimes) String() string {
  return fmt.Sprintf( "%8v %5v %5v %8v %8v %8v", &l.Time, l.Pta, l.Ptd, l.Wta, l.Wtd, l.Wtp )
}
