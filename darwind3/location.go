package darwind3

import (
  "encoding/xml"
  "fmt"
  "github.com/peter-mount/golib/codec"
  "strconv"
)

// A location in a schedule.
// This is formed of the entries from a schedule and is updated by any incoming
// Forecasts.
//
// As schedules can be circular (i.e. start and end at the same station) then
// the unique key is Tiploc and CircularTimes.Time.
//
// Location's within a Schedule are sorted by CircularTimes.Time accounting for
// crossing over midnight.
type Location struct {
  // Type of location, OR OPOR IP OPIP PP DT or OPDT
  Type              string        `json:"type"`
  // Tiploc of this location
  Tiploc            string        `json:"tiploc"`
  // The times for this entry
  Times             CircularTimes `json:"timetable"`
  // TIPLOC of False Destination to be used at this location
  FalseDestination  string        `json:"FalseDestination,omitempty"`
  // Is this service cancelled at this location
  Cancelled         bool          `json:"cancelled,omitempty"`
  // The Planned data for this location
  // i.e. information planned in advance
  Planned struct {
    // Current Activity Codes
    ActivityType      string      `json:"activity,omitempty"`
    // Planned Activity Codes (if different to current activities)
    PlannedActivity   string      `json:"plannedActivity,omitempty"`
    // A delay value that is implied by a change to the service's route.
    // This value has been added to the forecast lateness of the service at
    // the previous schedule location when calculating the expected lateness
    // of arrival at this location.
    RDelay            int         `json:"rDelay,omitempty"`
  }                               `json:"planned"`
  // The Forecast data at this location
  // i.e. information that changes in real time
  Forecast struct {
    // Forecast data for the arrival at this location
    Arrival           TSTime      `json:"arr"`
    // Forecast data for the departure at this location
    Departure         TSTime      `json:"dep"`
    // Forecast data for the pass of this location
    Pass              TSTime      `json:"pass"`
    // Current platform number
    Platform          Platform    `json:"plat"`
    // The service is suppressed at this location.
    Suppressed        bool        `json:"suppressed,omitempty"`
    // The length of the service at this location on departure
    // (or arrival at destination).
    // The default value of zero indicates that the length is unknown.
    Length            int         `json:"length,omitempty"`
    // Indicates from which end of the train stock will be detached.
    // The value is set to “true” if stock will be detached from the front of
    // the train at this location. It will be set at each location where stock
    // will be detached from the front.
    // Darwin will not validate that a stock detachment activity code applies
    // at this location.
    DetachFront       bool        `json:"detachFront,omitempty"`
    // The train order at this location (1, 2 or 3). 0 Means no TrainOrder has been set
    TrainOrder       *TrainOrder  `json:"trainOrder,omitempty"`
  }                               `json:"forecast"`
}

// Compare compares two Locations by their times
func (a *Location) Compare( b *Location ) bool {
  return b != nil && a.Times.Compare( &b.Times )
}

// Equals compares two Locations based on their Tiploc & time.
// This is used when trying to locate a location that's been updated
func (a *Location) EqualInSchedule( b *Location ) bool {
  return b != nil &&
         a.Tiploc == b.Tiploc &&
         a.Times.Time.Equals( &b.Times.Time )
}

// Equals compares two Locations in their entirety
func (a *Location) Equals( b *Location ) bool {
  return b != nil &&
         a.Type == b.Type &&
         a.Tiploc == b.Tiploc &&
         a.Times.Equals( &b.Times ) &&
         a.FalseDestination == b.FalseDestination &&
         a.Cancelled == b.Cancelled &&
         a.Planned.ActivityType == b.Planned.ActivityType &&
         a.Planned.PlannedActivity == b.Planned.PlannedActivity &&
         a.Planned.RDelay == b.Planned.RDelay &&
         a.Forecast.Arrival.Equals( &b.Forecast.Arrival ) &&
         a.Forecast.Departure.Equals( &b.Forecast.Departure ) &&
         a.Forecast.Pass.Equals( &b.Forecast.Pass ) &&
         a.Forecast.Platform.Equals( &b.Forecast.Platform ) &&
         a.Forecast.Length == b.Forecast.Length &&
         a.Forecast.DetachFront == b.Forecast.DetachFront &&
         a.Forecast.TrainOrder == b.Forecast.TrainOrder
}

func (t *Location) Write( c *codec.BinaryCodec ) {
  c.WriteString( t.Type ).
    WriteString( t.Tiploc ).
    WriteString( t.FalseDestination ).
    WriteBool( t.Cancelled )

  // CircularTimes
  c.Write( &t.Times )

  // Planned
  c.WriteString( t.Planned.ActivityType ).
    WriteString( t.Planned.PlannedActivity ).
    WriteInt( t.Planned.RDelay )

  // Forecast
  c.Write( &t.Forecast.Arrival ).
    Write( &t.Forecast.Departure ).
    Write( &t.Forecast.Pass ).
    Write( &t.Forecast.Platform ).
    WriteBool( t.Forecast.Suppressed ).
    WriteBool( t.Forecast.DetachFront ).
    WriteInt16( int16( t.Forecast.Length ) )

  if t.Forecast.TrainOrder == nil || t.Forecast.TrainOrder.Order == 0 {
    c.WriteByte( byte( 0 ) )
  } else {
    c.WriteByte( byte( t.Forecast.TrainOrder.Order ) ).
      WriteString( t.Forecast.TrainOrder.Platform )
  }
}

func (t *Location) Read( c *codec.BinaryCodec ) {
  c.ReadString( &t.Type ).
    ReadString( &t.Tiploc ).
    ReadString( &t.FalseDestination ).
    ReadBool( &t.Cancelled )

  // CircularTimes
  c.Read( &t.Times )

  // Planned
  c.ReadString( &t.Planned.ActivityType ).
    ReadString( &t.Planned.PlannedActivity ).
    ReadInt( &t.Planned.RDelay )

  // Forecast
  c.Read( &t.Forecast.Arrival ).
    Read( &t.Forecast.Departure ).
    Read( &t.Forecast.Pass ).
    Read( &t.Forecast.Platform ).
    ReadBool( &t.Forecast.Suppressed ).
    ReadBool( &t.Forecast.DetachFront )

  var l int16
  c.ReadInt16( &l )
  t.Forecast.Length = int(l)

  var b byte
  c.ReadByte( &b )
  if b > 0 {
    t.Forecast.TrainOrder = &TrainOrder{ Order: int( b ) }
    c.ReadString( &t.Forecast.TrainOrder.Platform )
  }
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

  // Parse CircularTimes attributes
  s.Times.UnmarshalXMLAttributes( start )

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
          if err := decoder.DecodeElement( elem, &tok ); err != nil {
            return err
          }
        }

      case xml.EndElement:
        return nil
    }
  }
}

func (l *Location) String() string {
  return fmt.Sprintf(
    "%2s %7s %7s %5v %s %s %s %d %v %v %v %v %v %d %v",
    l.Type,
    l.Tiploc,
    l.FalseDestination,
    l.Cancelled,
    l.Times.String(),
    l.Planned.ActivityType,
    l.Planned.PlannedActivity,
    l.Planned.RDelay,
    l.Forecast.Arrival,
    l.Forecast.Departure,
    l.Forecast.Pass,
    l.Forecast.Platform,
    l.Forecast.Suppressed,
    l.Forecast.Length,
    l.Forecast.DetachFront )
}
