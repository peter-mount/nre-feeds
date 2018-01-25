package darwind3

import (
  "github.com/peter-mount/golib/codec"
)

// Platform number with associated flags
type Platform struct {
  // Defines a platform number
  Platform          string    `xml:",chardata"`
  // True if the platform number is confirmed.
  Confirmed         bool      `xml:"conf,attr,omitempty"`
  // Platform number is suppressed and should not be displayed.
  Suppressed        bool      `xml:"platsup,attr,omitempty"`
  // Whether a CIS, or Darwin Workstation, has set platform suppression at this location.
  CISSuppressed     bool      `xml:"cisPlatsup,attr,omitempty"`
  // The source of the platfom number. P = Planned, A = Automatic, M = Manual.
  // Default is P
  Source            string    `xml:"platsrc,attr,omitempty"`
}

func (a *Platform ) Equals( b *Platform ) bool {
  return b != nil &&
         a.Platform == b.Platform &&
         a.Confirmed == b.Confirmed &&
         a.Suppressed == b.Suppressed &&
         a.CISSuppressed == b.CISSuppressed &&
         a.Source == b.Source
}

func (t *Platform) Write( c *codec.BinaryCodec ) {
  c.WriteString( t.Platform ).
    WriteBool( t.Confirmed ).
    WriteBool( t.Suppressed ).
    WriteBool( t.CISSuppressed ).
    WriteString( t.Source )
}

func (t *Platform) Read( c *codec.BinaryCodec ) {
  c.ReadString( &t.Platform ).
    ReadBool( &t.Confirmed ).
    ReadBool( &t.Suppressed ).
    ReadBool( &t.CISSuppressed ).
    ReadString( &t.Source )
}
