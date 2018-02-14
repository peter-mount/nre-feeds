package darwind3

import (
  "bytes"
  "darwintimetable"
  "encoding/json"
  "encoding/xml"
  "github.com/peter-mount/golib/codec"
)

// Type describing time-based forecast attributes for a TS arrival/departure/pass
type TSTime struct {
  // Estimated Time. For locations with a public activity,
  // this will be based on the "public schedule".
  // For all other activities, it will be based on the "working schedule".
  ET         *darwintimetable.WorkingTime   `json:"et,omitempty" xml:"et,attr,omitempty"`
  // The manually applied lower limit that has been applied to the estimated
  // time at this location. The estimated time will not be set lower than this
  // value, but may be set higher.
  ETMin      *darwintimetable.WorkingTime   `json:"etMin,omitempty" xml:"etmin,attr,omitempty"`
  // Indicates that an unknown delay forecast has been set for the estimated
  // time at this location. Note that this value indicates where a manual
  // unknown delay forecast has been set, whereas it is the "delayed"
  // attribute that indicates that the actual forecast is "unknown delay".
  ETUnknown   bool                          `json:"etUnknown,omitempty" xml:"etUnknown,attr,omitempty"`
  // The estimated time based on the "working schedule".
  // This will only be set for public activities and when it also differs
  // from the estimated time based on the "public schedule".
  WET        *darwintimetable.WorkingTime   `json:"wet,omitempty" xml:"wet,attr,omitempty"`
  // Actual Time
  AT         *darwintimetable.WorkingTime   `json:"at,omitempty" xml:"at,attr,omitempty"`
  // If true, indicates that an actual time ("at") value has just been removed
  // and replaced by an estimated time ("et").
  // Note that this attribute will only be set to "true" once, when the actual
  // time is removed, and will not be set in any snapshot.
  ATRemoved   bool        `json:"atRemoved,omitempty" xml:"atRemoved,attr,omitempty"`
  // Indicates that this estimated time is a forecast of "unknown delay".
  // Displayed  as "Delayed" in LDB.
  // Note that this value indicates that this forecast is "unknown delay",
  // whereas it is the "etUnknown" attribute that indicates where the manual
  // unknown delay forecast has been set.
  Delayed     bool        `json:"delayed,omitempty" xml:"delayed,attr,omitempty"`
  // The source of the forecast or actual time.
  Src         string      `json:"src,omitempty" xml:"src,attr,omitempty"`
  // The RTTI CIS code of the CIS instance if the src is a CIS.
  SrcInst     string      `json:"srcInst,omitempty" xml:"srcInst,attr,omitempty"`
}

func (t *TSTime) append( b *bytes.Buffer, c bool, f string, v interface{} ) bool {
  // Any null, "" or false ignore
  if vb, err := json.Marshal( v );
    err == nil &&
    !( len(vb) == 2 && vb[0] == '"' && vb[1] == '"' ) &&
    !( len(vb) == 4 && vb[0] == 'n' && vb[1] == 'u' && vb[2] == 'l' && vb[3] == 'l') &&
    !( len(vb) == 5 && vb[0] == 'f' && vb[1] == 'a' && vb[2] == 'l' && vb[3] == 's' && vb[4] == 'e') {
    if c {
      b.WriteByte( ',' )
    }

    b.WriteByte( '"' )
    b.WriteString( f )
    b.WriteByte( '"' )
    b.WriteByte( ':' )
    b.Write( vb )
    return true
  }

  return c
}

func (t *TSTime) MarshalJSON() ( []byte, error ) {
  var b bytes.Buffer

  if t.ET == nil && t.ETMin == nil && t.AT == nil && t.WET == nil {
    b.WriteString( "null" )
  } else {
    b.WriteByte( '{' )
    c := t.append( &b, false, "et", t.ET )
    c = t.append( &b, c, "etMin", t.ETMin )
    c = t.append( &b, c, "etMin", t.ETMin )
    c = t.append( &b, c, "etUnknown", &t.ETUnknown )
    c = t.append( &b, c, "wet", t.WET )
    c = t.append( &b, c, "at", t.AT )
    c = t.append( &b, c, "atRemoved", &t.ATRemoved )
    c = t.append( &b, c, "delayed", &t.Delayed )
    c = t.append( &b, c, "src", &t.Src )
    c = t.append( &b, c, "srcInst", &t.SrcInst )
    b.WriteByte( '}' )
  }

  return b.Bytes(), nil
}

// IsSet returns true if at least ET or AT is set
func (t *TSTime) IsSet() bool {
  return ( t.ET != nil && !t.ET.IsZero() ) || ( t.AT !=nil && !t.AT.IsZero() )
}

// Compare compares two TSTime's.
// This will use the value returned by TSTime.Time()
func (a *TSTime) Compare( b *TSTime ) bool {
  if b == nil {
    return false
  }
  var at = a.Time()
  var bt = b.Time()
  return at != nil && at.Compare( bt )
}

// Time returns the appropirate time from TSTime to use in displays.
// This is the first one set of AT, ET or nil if neither is set.
func (t *TSTime) Time() *darwintimetable.WorkingTime {
  if t.AT != nil {
    return t.AT
  }
  if t.ET != nil {
    return t.ET
  }
  return nil
}

func (a *TSTime) Equals( b *TSTime ) bool {
  return b != nil &&
         a.ET == b.ET &&
         a.ETMin == b.ETMin &&
         a.ETUnknown == b.ETUnknown &&
         a.WET == b.WET &&
         a.AT == b.AT &&
         a.ATRemoved == b.ATRemoved &&
         a.Delayed == b.Delayed &&
         a.Src == b.Src &&
         a.SrcInst == b.SrcInst
}

func (t *TSTime) Write( c *codec.BinaryCodec ) {
  darwintimetable.WorkingTimeWrite( c, t.ET )
  darwintimetable.WorkingTimeWrite( c, t.ETMin )
  darwintimetable.WorkingTimeWrite( c, t.WET )
  darwintimetable.WorkingTimeWrite( c, t.AT )
  c.WriteBool( t.ETUnknown ).
    WriteBool( t.ATRemoved ).
    WriteBool( t.Delayed ).
    WriteString( t.Src ).
    WriteString( t.SrcInst )
}

func (t *TSTime) Read( c *codec.BinaryCodec ) {
  t.ET = darwintimetable.WorkingTimeRead( c )
  t.ETMin = darwintimetable.WorkingTimeRead( c )
  t.WET = darwintimetable.WorkingTimeRead( c )
  t.AT = darwintimetable.WorkingTimeRead( c )
  c.ReadBool( &t.ETUnknown ).
    ReadBool( &t.ATRemoved ).
    ReadBool( &t.Delayed ).
    ReadString( &t.Src ).
    ReadString( &t.SrcInst )
}

func (s *TSTime) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {

  for _, attr := range start.Attr {
    switch attr.Name.Local {
      case "et":
        s.ET = darwintimetable.NewWorkingTime( attr.Value )

      case "etmin":
        s.ETMin = darwintimetable.NewWorkingTime( attr.Value )

      case "etUnknown":
        s.ETUnknown = attr.Value == "true"

      case "wet":
        s.WET = darwintimetable.NewWorkingTime( attr.Value )

      case "at":
        s.AT = darwintimetable.NewWorkingTime( attr.Value )

      case "atRemoved":
        s.ATRemoved = attr.Value == "true"

      case "delayed":
        s.Delayed = attr.Value == "true"

      case "src":
        s.Src = attr.Value

      case "srcInst":
        s.SrcInst = attr.Value

    }
  }

  return decoder.Skip()
}
