package darwind3

import (
  "darwintimetable"
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
        s.SrcInst = attr.Value

      case "srcInst":
        s.SrcInst = attr.Value

    }
  }

  return decoder.Skip()
}
