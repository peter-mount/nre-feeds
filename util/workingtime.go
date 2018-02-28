package util

import (
  "encoding/json"
  "encoding/xml"
  "fmt"
  "github.com/peter-mount/golib/codec"
  "strconv"
  "time"
)

// Working Timetable time.
// WorkingTime is similar to PublciTime, except we can have seconds.
// In the Working Timetable, the seconds can be either 0 or 30.
type WorkingTime struct {
  t int
}

const (
  workingTime_min = 0
  workingTime_max = 86400
)

func (a *WorkingTime) Equals( b *WorkingTime ) bool {
  if b == nil {
    return false
  }
  if a.IsZero() {
    return b.IsZero()
  }
  return a.t == b.t
}

// Compare a WorkingTime against another, accounting for crossing midnight.
// The rules for handling crossing midnight are:
// < -6 hours = crossed midnight
// < 0 back in time
// < 18 hours increasing time
// > 18 hours back in time & crossing midnight
func (a *WorkingTime) Compare( b *WorkingTime ) bool {
  if b == nil {
    return false
  }

  d := a.t - b.t

  if d < -21600 || d >64800 {
    return a.t > b.t
  }

  return a.t < b.t
}

// WorkingTimeEquals compares equality between two WorkingTimes.
// Unlike WorkingTime.Equals() this will return true if both are null, otherwise
// both must not be null and equal to be true
func WorkingTimeEquals( a *WorkingTime, b *WorkingTime ) bool {
  return ( a == nil && b == nil ) ||
         ( a != nil && a.Equals( b ) )
}

// BinaryCodec writer
func (t *WorkingTime) Write( c *codec.BinaryCodec ) {
  if t.IsZero() {
    // If zero then write a single byte value 255
    c.WriteByte( byte(255) )
  } else {
    // 24 bit number, Big endian so we should not hit 255,xxx,xxx in the stream
    c.WriteByte( byte( t.t >> 16 ) ).
      WriteByte( byte( t.t >> 8 ) ).
      WriteByte( byte( t.t ) )
  }
}

// BinaryCodec reader
func (t *WorkingTime) Read( c *codec.BinaryCodec ) {
  var b byte

  c.ReadByte( &b )
  if b == byte(255) {
    t.t = -1
  } else {
    var i uint64 = uint64(b)<<16
    c.ReadByte( &b )
    i |= uint64(b)<<8
    c.ReadByte( &b )
    t.t = int( i | uint64(b) )
  }
}

// NewWorkingTime returns a new WorkingTime instance from a string of format "HH:MM:SS"
func NewWorkingTime( s string ) *WorkingTime {
  v := &WorkingTime{}
  v.Parse( s )
  return v
}

func (v *WorkingTime) Parse( s string ) {
  if s == "" {
    v.t = -1
  } else {
    a, _ := strconv.Atoi( s[0:2] )
    b, _ := strconv.Atoi( s[3:5] )
    if len( s ) > 5 {
      c, _ := strconv.Atoi( s[6:8] )
      v.Set( (a *3600) + (b * 60) + c )
    } else {
      v.Set( (a *3600) + (b * 60) )
    }
  }
}

// WorkingTimeRead is a workaround issue where a custom type cannot be
// omitempty in JSON unless it's a nil
// So instead of using BinaryCodec.Read( v ), we call this & set the return
// value in the struct as a pointer.
func WorkingTimeRead( c *codec.BinaryCodec ) *WorkingTime {
  t := &WorkingTime{}
  c.Read( t )
  if t.IsZero() {
    return nil
  }
  return t
}

// WorkingTimeWrite is a workaround for writing null times.
// If the pointer is null then a time is written where IsZero()==true
func WorkingTimeWrite( c *codec.BinaryCodec, t *WorkingTime ) {
  if t == nil {
    // Null so mark as zero
    c.WriteByte( byte(255) )
  } else {
    c.Write( t )
  }
}

// Custom JSON Marshaler. This will write null or the time as "HH:MM:SS"
func (t *WorkingTime) MarshalJSON() ( []byte, error ) {
  if t.IsZero() {
    return json.Marshal( nil )
  }
  return json.Marshal( t.String() )
}

func (t *WorkingTime) UnmarshalJSON( b []byte ) error {
  s := string(b[:])
  if s != "null" && len( s ) > 2 {
    t.Parse( s[1:len(s)-1] )
  }
  return nil
}

// Custom XML Marshaler.
func (t *WorkingTime) MarshalXMLAttr( name xml.Name ) ( xml.Attr, error ) {
  if t.IsZero() {
    return xml.Attr{}, nil
  }
  return xml.Attr{ Name: name, Value: t.String() }, nil
}

// String returns a PublicTime in HH:MM:SS format or 8 blank spaces if it's not set.
func (t *WorkingTime) String() string {
  if t.IsZero() {
    return "        "
  }

  return fmt.Sprintf( "%02d:%02d:%02d", t.t/3600, (t.t/60)%60, t.t%60 )
}

// Get returns the WorkingTime in seconds of the day
func (t *WorkingTime) Get() int {
  return t.t
}

// Set sets the WorkingTime in seconds of the day
func (t *WorkingTime) Set( v int ) {
  t.t = v
}

// IsZero returns true if the time is not present
func (t *WorkingTime) IsZero() bool {
  return t.t <= 0
}

// SetTime set's the working time to the current time (resolution 1 minute)
func (t *WorkingTime) SetTime( tm time.Time ) {
  t.Set( (tm.Hour()*3600) + (tm.Minute()*60) )
}

// WorkingTime_FromTime returns a WorkingTime from a time.Time with a resolution
// of 1 minute.
func WorkingTime_FromTime( tm time.Time ) *WorkingTime {
  t := &WorkingTime{}
  t.SetTime( tm )
  return t
}

// Before returns true if this WorkingTime is before another
func (a *WorkingTime) Before( b *WorkingTime ) bool {
  return a.t < b.t
}

// After returns true if this WorkingTime is after another
func (a *WorkingTime) After( b *WorkingTime ) bool {
  return a.t > b.t
}

// Between returns true if this WorkingTime falls between two other WorkingTime's.
// The test is inclusive of the from & to times.
// If from is after to then we presume we cross midnight.
func (t *WorkingTime) Between( from *WorkingTime, to *WorkingTime ) bool {

  if from.After( to ) {
    return ( from.t <= t.t && t.t <= workingTime_max ) || ( t.t >= workingTime_min && t.t <= to.t )
  }

  return from.t <= t.t && t.t <= to.t
}
