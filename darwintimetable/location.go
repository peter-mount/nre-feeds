// Reference timetable
package darwintimetable

import (
  //bolt "github.com/coreos/bbolt"
  "encoding/xml"
  "github.com/peter-mount/golib/codec"
)

// Common location object used in persistence
type Location struct {
    XMLName     xml.Name    `json:"-" xml:"location"`
    Type        string      `json:"type" xml:"type,attr"`
    Tiploc      string      `json:"tpl" xml:"tpl,attr"`
    Act         string      `json:"act,omitempty" xml:"act,attr,omitempty"`
    PlanAct     string      `json:"planAct,omitempty" xml:"planAct,attr,omitempty"`
    Cancelled   bool        `json:"cancelled,omitempty" xml:"can,attr,omitempty"`
    Platform    string      `json:"plat,omitempty" xml:"plat,attr,omitempty"`
    // CallPtAttributes
    Pta        *PublicTime  `json:"pta,omitempty" xml:"pta,attr,omitempty"`
    Ptd        *PublicTime  `json:"ptd,omitempty" xml:"ptd,attr,omitempty"`
    // Working times
    Wta        *WorkingTime `json:"wta,omitempty" xml:"wta,attr,omitempty"`
    Wtd        *WorkingTime `json:"wtd,omitempty" xml:"wtd,attr,omitempty"`
    Wtp        *WorkingTime `json:"wtp,omitempty" xml:"wtp,attr,omitempty"`
    // Delay implied by a change to the services route
    RDelay      string      `json:"rdelay,omitempty" xml:"rdelay,attr,omitempty"`
    // False destination to be used at this location
    FalseDest   string      `json:"fd,omitempty" xml:"fd,attr,omitempty"`
}

func (t *Location) Write( c *codec.BinaryCodec ) {
  c.WriteString( t.Type ).
    WriteString( t.Tiploc ).
    WriteString( t.Act ).
    WriteString( t.PlanAct ).
    WriteBool( t.Cancelled ).
    WriteString( t.Platform ).
    WriteString( t.RDelay ).
    WriteString( t.FalseDest )
  PublicTimeWrite( c, t.Pta )
  PublicTimeWrite( c, t.Ptd )
  WorkingTimeWrite( c, t.Wta )
  WorkingTimeWrite( c, t.Wtd )
  WorkingTimeWrite( c, t.Wtp )
}

func (t *Location) Read( c *codec.BinaryCodec ) {
  c.ReadString( &t.Type ).
    ReadString( &t.Tiploc ).
    ReadString( &t.Act ).
    ReadString( &t.PlanAct ).
    ReadBool( &t.Cancelled ).
    ReadString( &t.Platform ).
    ReadString( &t.RDelay ).
    ReadString( &t.FalseDest )
  t.Pta = PublicTimeRead( c )
  t.Ptd = PublicTimeRead( c )
  t.Wta = WorkingTimeRead( c )
  t.Wtd = WorkingTimeRead( c )
  t.Wtp = WorkingTimeRead( c )
}

type OR struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // CallPtAttributes
  Pta       string      `xml:"pta,attr"`
  Ptd       string      `xml:"ptd,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // False destination to be used at this location
  FalseDest string      `xml:"fd,attr"`
}

func (s *OR ) Location() *Location {
  return &Location{
    Type: "OR",
    Tiploc: s.Tiploc,
    Act: s.Act,
    PlanAct: s.PlanAct,
    Cancelled: s.Cancelled,
    Platform: s.Platform,
    Pta: NewPublicTime( s.Pta ),
    Ptd: NewPublicTime( s.Ptd ),
    Wta: NewWorkingTime( s.Wta ),
    Wtd: NewWorkingTime( s.Wtd ),
    FalseDest: s.FalseDest,
  }
}

type OPOR struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
}

func (s *OPOR ) Location() *Location {
  return &Location{
    Type: "OPOR",
    Tiploc: s.Tiploc,
    Act: s.Act,
    PlanAct: s.PlanAct,
    Cancelled: s.Cancelled,
    Platform: s.Platform,
    Wta: NewWorkingTime( s.Wta ),
    Wtd: NewWorkingTime( s.Wtd ),
  }
}

type IP struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // CallPtAttributes
  Pta       string      `xml:"pta,attr"`
  Ptd       string      `xml:"ptd,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
  // False destination to be used at this location
  FalseDest string      `xml:"fd,attr"`
}

func (s *IP ) Location() *Location {
  return &Location{
    Type: "IP",
    Tiploc: s.Tiploc,
    Act: s.Act,
    PlanAct: s.PlanAct,
    Cancelled: s.Cancelled,
    Platform: s.Platform,
    Pta: NewPublicTime( s.Pta ),
    Ptd: NewPublicTime( s.Ptd ),
    Wta: NewWorkingTime( s.Wta ),
    Wtd: NewWorkingTime( s.Wtd ),
    RDelay: s.RDelay,
    FalseDest: s.FalseDest,
  }
}

type OPIP struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

func (s *OPIP ) Location() *Location {
  return &Location{
    Type: "OPIP",
    Tiploc: s.Tiploc,
    Act: s.Act,
    PlanAct: s.PlanAct,
    Cancelled: s.Cancelled,
    Platform: s.Platform,
    Wta: NewWorkingTime( s.Wta ),
    Wtd: NewWorkingTime( s.Wtd ),
    RDelay: s.RDelay,
  }
}

type PP struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wtp       string      `xml:"wtp,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

func (s *PP ) Location() *Location {
  return &Location{
    Type: "PP",
    Tiploc: s.Tiploc,
    Act: s.Act,
    PlanAct: s.PlanAct,
    Cancelled: s.Cancelled,
    Platform: s.Platform,
    Wtp: NewWorkingTime( s.Wtp ),
    RDelay: s.RDelay,
  }
}

type DT struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // CallPtAttributes
  Pta       string      `xml:"pta,attr"`
  Ptd       string      `xml:"ptd,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

func (s *DT ) Location() *Location {
  return &Location{
    Type: "DT",
    Tiploc: s.Tiploc,
    Act: s.Act,
    PlanAct: s.PlanAct,
    Cancelled: s.Cancelled,
    Platform: s.Platform,
    Pta: NewPublicTime( s.Pta ),
    Ptd: NewPublicTime( s.Ptd ),
    Wta: NewWorkingTime( s.Wta ),
    Wtd: NewWorkingTime( s.Wtd ),
    RDelay: s.RDelay,
  }
}

type OPDT struct {
  // SchedLocAttributes
  Tiploc    string      `xml:"tpl,attr"`
  Act       string      `xml:"act,attr"`
  PlanAct   string      `xml:"planAct,attr"`
  Cancelled bool        `xml:"can,attr"`
  Platform  string      `xml:"plat,attr"`
  // Working times
  Wta       string      `xml:"wta,attr"`
  Wtd       string      `xml:"wtd,attr"`
  // Delay implied by a change to the services route
  RDelay    string      `xml:"rdelay,attr"`
}

func (s *OPDT ) Location() *Location {
  return &Location{
    Type: "OPDT",
    Tiploc: s.Tiploc,
    Act: s.Act,
    PlanAct: s.PlanAct,
    Cancelled: s.Cancelled,
    Platform: s.Platform,
    Wta: NewWorkingTime( s.Wta ),
    Wtd: NewWorkingTime( s.Wtd ),
    RDelay: s.RDelay,
  }
}
