package darwind3

import (
  "github.com/peter-mount/golib/codec"
)

// Type used to represent a cancellation or late running reason
type DisruptionReason struct {
  // A Darwin Reason Code. 0 = none
  Reason  int       `json:"reason" xml:",chardata"`
  // Optional TIPLOC where the reason refers to, e.g. "signalling failure at Cheadle Hulme"
  Tiploc  string    `json:"tiploc,omitempty" xml:"tiploc,attr,omitempty"`
  // If true, the tiploc attribute should be interpreted as "near",
  // e.g. "signalling failure near Cheadle Hulme".
  Near    bool      `json:"near,omitempty" xml:"near,attr,omitempty"`
}

func (a *DisruptionReason) Equals( b *DisruptionReason ) bool {
  return b != nil &&
         a.Reason == b.Reason &&
         a.Tiploc == b.Tiploc &&
         a.Near == b.Near
}

func (t *DisruptionReason) Write( c *codec.BinaryCodec ) {
  c.WriteInt16( int16( t.Reason ) ).
    WriteString( t.Tiploc ).
    WriteBool( t.Near )
}

func (t *DisruptionReason) Read( c *codec.BinaryCodec ) {
  var r int16
  c.ReadInt16( &r )
  t.Reason = int(r)

  c.ReadString( &t.Tiploc ).
    ReadBool( &t.Near )
}
