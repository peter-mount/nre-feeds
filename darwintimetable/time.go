package darwintimetable

import (
  "encoding/json"
  "encoding/xml"
  "fmt"
  "github.com/peter-mount/golib/codec"
  "strconv"
)

// Public Timetable time
// Note: 00:00 is not possible as in CIF that means no-time
type PublicTime struct {
  t int
}

// BinaryCodec writer
func (t *PublicTime) Write( c *codec.BinaryCodec ) {
  c.WriteInt32( int32( t.t ) )
}

// BinaryCodec reader
func (t *PublicTime) Read( c *codec.BinaryCodec ) {
  var i int32
  c.ReadInt32( &i )
  t.t = int(i)
}

func NewPublicTime( s string ) *PublicTime {
  v := &PublicTime{}
  if s == "" {
    v.t = -1
  } else {
    a, _ := strconv.Atoi( s[0:2] )
    b, _ := strconv.Atoi( s[2:4] )
    v.Set( (a *3600) + (b * 60) )
  }
  return v
}

// PublicTimeRead is a workaround issue where a custom type cannot be
// omitempty in JSON unless it's a nil
// So instead of using BinaryCodec.Read( v ), we call this & set the return
// value in the struct as a pointer.
func PublicTimeRead( c *codec.BinaryCodec ) *PublicTime {
  t := &PublicTime{}
  c.Read( t )
  if t.IsZero() {
    return nil
  }
  return t
}

// PublicTimeWrite is a workaround for writing null times.
// If the pointer is null then a time is written where IsZero()==true
func PublicTimeWrite( c *codec.BinaryCodec, t *PublicTime ) {
  if t == nil {
    c.WriteInt32( -1 )
  } else {
    c.Write( t )
  }
}

// Custom JSON Marshaler. This will write null or the time as "HH:MM"
func (t *PublicTime) MarshalJSON() ( []byte, error ) {
  if t.IsZero() {
    return json.Marshal( nil )
  }
  return json.Marshal( t.String() )
}

// Custom XML Marshaler.
func (t *PublicTime) MarshalXMLAttr( name xml.Name ) ( xml.Attr, error ) {
  if t.IsZero() {
    return xml.Attr{}, nil
  }
  return xml.Attr{ Name: name, Value: t.String() }, nil
}

// String returns a PublicTime in HH:MM format or 5 blank spaces if it's not set.
func (t *PublicTime) String() string {
  if t.t <= 0 {
    return "     "
  }

  return fmt.Sprintf( "%02d:%02d", t.t/3600, (t.t/60)%60 )
}

// Get returns the PublicTime in seconds of the day
func (t *PublicTime) Get() int {
  return t.t
}

// Set sets the PublicTime in seconds of the day
func (t *PublicTime) Set( v int ) {
  t.t = v
}

// IsZero returns true if the time is not present
func (t *PublicTime) IsZero() bool {
  return t.t <= 0
}

// Working Timetable time.
// WorkingTime is similar to PublciTime, except we can have seconds.
// In the Working Timetable, the seconds can be either 0 or 30.
type WorkingTime struct {
  t int
}

// BinaryCodec writer
func (t *WorkingTime) Write( c *codec.BinaryCodec ) {
  c.WriteInt32( int32( t.t ) )
}

// BinaryCodec reader
func (t *WorkingTime) Read( c *codec.BinaryCodec ) {
  var i int32
  c.ReadInt32( &i )
  t.t = int(i)
}

func NewWorkingTime( s string ) *WorkingTime {
  v := &WorkingTime{}
  if s == "" {
    v.t = -1
  } else {
    a, _ := strconv.Atoi( s[0:2] )
    b, _ := strconv.Atoi( s[3:5] )
    if len( s ) > 6 {
      c, _ := strconv.Atoi( s[6:8] )
      v.Set( (a *3600) + (b * 60) + c )
    } else {
      v.Set( (a *3600) + (b * 60) )
    }
  }
  return v
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
    c.WriteInt32( -1 )
  } else {
    c.Write( t )
  }
}

// Custom JSON Marshaler. This will write null or the time as "HH:MM:SS"
func (t *WorkingTime) MarshalJSON() ( []byte, error ) {
  if t.t < 0 {
    return json.Marshal( nil )
  }
  return json.Marshal( t.String() )
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
  return t.t < 0
}
