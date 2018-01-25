package darwind3

import (
  "darwintimetable"
  "encoding/xml"
  "fmt"
  "github.com/peter-mount/golib/codec"
  "strconv"
)

type Location struct {
  // Type of location, OR OPOR IP OPIP PP DT or OPDT
  Type              string
  // Tiploc of this location
  Tiploc            string
  // TIPLOC of False Destination to be used at this location
  FalseDestination  string
  // Is this service cancelled at this location
  Cancelled         bool
  // The time for this location
  Time                darwintimetable.WorkingTime
  // The scheduled data for this location
  Planned struct {
    // Current Activity Codes
    ActivityType      string
    // Planned Activity Codes (if different to current activities)
    PlannedActivity   string
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
    // A delay value that is implied by a change to the service's route.
    // This value has been added to the forecast lateness of the service at
    // the previous schedule location when calculating the expected lateness
    // of arrival at this location.
    RDelay            int
  }
  // The Forecast data at this location
  Forecast struct {
    // Forecast data for the arrival at this location
    Arrival           TSTime
    // Forecast data for the departure at this location
    Departure         TSTime
    // Forecast data for the pass of this location
    Pass              TSTime
    // Current platform number
    Platform          Platform
    // The service is suppressed at this location.
    Suppressed        bool
    // The length of the service at this location on departure
    // (or arrival at destination).
    // The default value of zero indicates that the length is unknown.
    Length            int
    // Indicates from which end of the train stock will be detached.
    // The value is set to “true” if stock will be detached from the front of
    // the train at this location. It will be set at each location where stock
    // will be detached from the front.
    // Darwin will not validate that a stock detachment activity code applies
    // at this location.
    DetachFront       bool
  }
}

// UpdateTime updates the Time field used for sequencing the location.
// This is the the first one of these set in the following order:
// Wtd, Wta, Wtp
// Note this value is not persisted as it's a generated value
func (l *Location) UpdateTime() {
  if l.Planned.Wtd != nil {
    l.Time = *l.Planned.Wtd
  } else if l.Planned.Wta != nil {
    l.Time = *l.Planned.Wta
  } else {
    l.Time = *l.Planned.Wtp
  }
}

// Compare compares two Locations by their times
func (a *Location) Compare( b *Location ) bool {
  return b != nil && a.Time.Compare( &b.Time )
}

// Equals compares two Locations in their entirety
func (a *Location) Equals( b *Location ) bool {
  return b != nil &&
         a.Type == b.Type &&
         a.Tiploc == b.Tiploc &&
         a.FalseDestination == b.FalseDestination &&
         a.Cancelled == b.Cancelled &&
         a.Planned.ActivityType == b.Planned.ActivityType &&
         a.Planned.PlannedActivity == b.Planned.PlannedActivity &&
         darwintimetable.PublicTimeEquals( a.Planned.Pta, b.Planned.Pta ) &&
         darwintimetable.PublicTimeEquals( a.Planned.Ptd, b.Planned.Ptd ) &&
         darwintimetable.WorkingTimeEquals( a.Planned.Wta, b.Planned.Wta ) &&
         darwintimetable.WorkingTimeEquals( a.Planned.Wtd, b.Planned.Wtd ) &&
         darwintimetable.WorkingTimeEquals( a.Planned.Wtp, b.Planned.Wtp ) &&
         a.Planned.RDelay == b.Planned.RDelay &&
         a.Forecast.Arrival.Equals( &b.Forecast.Arrival ) &&
         a.Forecast.Departure.Equals( &b.Forecast.Departure ) &&
         a.Forecast.Pass.Equals( &b.Forecast.Pass ) &&
         a.Forecast.Platform.Equals( &b.Forecast.Platform ) &&
         a.Forecast.Length == b.Forecast.Length &&
         a.Forecast.DetachFront == b.Forecast.DetachFront
}

func (t *Location) Write( c *codec.BinaryCodec ) {
  c.WriteString( t.Type ).
    WriteString( t.Tiploc ).
    WriteString( t.FalseDestination ).
    WriteBool( t.Cancelled )

  // Planned
  c.WriteString( t.Planned.ActivityType ).
    WriteString( t.Planned.PlannedActivity ).
    WriteInt( t.Planned.RDelay )
  darwintimetable.PublicTimeWrite( c, t.Planned.Pta )
  darwintimetable.PublicTimeWrite( c, t.Planned.Ptd )
  darwintimetable.WorkingTimeWrite( c, t.Planned.Wta )
  darwintimetable.WorkingTimeWrite( c, t.Planned.Wtd )
  darwintimetable.WorkingTimeWrite( c, t.Planned.Wtp )

  // Forecast
  c.Write( &t.Forecast.Arrival ).
    Write( &t.Forecast.Departure ).
    Write( &t.Forecast.Pass ).
    Write( &t.Forecast.Platform ).
    WriteBool( t.Forecast.Suppressed ).
    WriteInt16( int16( t.Forecast.Length ) ).
    WriteBool( t.Forecast.DetachFront )
}

func (t *Location) Read( c *codec.BinaryCodec ) {
  c.ReadString( &t.Type ).
    ReadString( &t.Tiploc ).
    ReadString( &t.FalseDestination ).
    ReadBool( &t.Cancelled )

  // Planned
  c.ReadString( &t.Planned.ActivityType ).
    ReadString( &t.Planned.PlannedActivity ).
    ReadInt( &t.Planned.RDelay )
  t.Planned.Pta = darwintimetable.PublicTimeRead( c )
  t.Planned.Ptd = darwintimetable.PublicTimeRead( c )
  t.Planned.Wta = darwintimetable.WorkingTimeRead( c )
  t.Planned.Wtd = darwintimetable.WorkingTimeRead( c )
  t.Planned.Wtp = darwintimetable.WorkingTimeRead( c )

  // Forecast
  c.Read( &t.Forecast.Arrival ).
    Read( &t.Forecast.Departure ).
    Read( &t.Forecast.Pass ).
    Read( &t.Forecast.Platform ).
    ReadBool( &t.Forecast.Suppressed )

  var l int16
  c.ReadInt16( &l )
  t.Forecast.Length = int(l)

  c.ReadBool( &t.Forecast.DetachFront )

  // Update the time field
  t.UpdateTime()
}

func (s *Location) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {

  for _, attr := range start.Attr {
    switch attr.Name.Local {
      case "tpl":
        s.Tiploc = attr.Value

      case "act":
        s.Planned.ActivityType = attr.Value

      case "planAct":
        s.Planned.PlannedActivity = attr.Value

      case "Cancelled":
        s.Cancelled = attr.Value == "true"

      case "pta":
        s.Planned.Pta = darwintimetable.NewPublicTime( attr.Value )

      case "ptd":
        s.Planned.Ptd = darwintimetable.NewPublicTime( attr.Value )

      case "wta":
        s.Planned.Wta = darwintimetable.NewWorkingTime( attr.Value )

      case "wtd":
        s.Planned.Wtd = darwintimetable.NewWorkingTime( attr.Value )

      case "wtp":
        s.Planned.Wtp = darwintimetable.NewWorkingTime( attr.Value )

      case "fd":
        s.FalseDestination = attr.Value

      case "rdelay":
        if v, err := strconv.Atoi( attr.Value ); err != nil {
          return err
        } else {
          s.Planned.RDelay = v
        }
    }
  }

  // TODO parse body. Under schedule there is none but under TS there are

  // Update the time field
  s.UpdateTime()

  return decoder.Skip()
}

func (l *Location) String() string {
  return fmt.Sprintf(
    "%2s %7s %7s %5v %s %s %s %v %v %v %v %v %d %v %v %v %v %v %d %v",
    l.Type,
    l.Tiploc,
    l.FalseDestination,
    l.Cancelled,
    l.Time.String(),
    l.Planned.ActivityType,
    l.Planned.PlannedActivity,
    l.Planned.Pta,
    l.Planned.Ptd,
    l.Planned.Wta,
    l.Planned.Wtd,
    l.Planned.Wtp,
    l.Planned.RDelay,
    l.Forecast.Arrival,
    l.Forecast.Departure,
    l.Forecast.Pass,
    l.Forecast.Platform,
    l.Forecast.Suppressed,
    l.Forecast.Length,
    l.Forecast.DetachFront )
}
