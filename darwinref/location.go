package darwinref

import (
  "github.com/peter-mount/golib/codec"
  "time"
)
// Defines a location, i.e. Station or passing point
type Location struct {
  Tiploc      string    `json:"tpl" xml:"tpl,attr"`
  Crs         string    `json:"crs,omitempty" xml:"crs,attr,omitempty"`
  Toc         string    `json:"toc,omitempty" xml:"toc,attr,omitempty"`
  Name        string    `json:"locname" xml:"locname,attr"`
  // The CIF extract this entry is from
  Date        time.Time `json:"date" xml:"date,attr"`
}

func (a *Location) Equals( b *Location ) bool {
  if b == nil {
    return false
  }
  return a.Tiploc == b.Tiploc &&
    a.Crs == b.Crs &&
    a.Toc == b.Toc &&
    a.Name == b.Name
}

func (t *Location) Write( c *codec.BinaryCodec ) {
  c.WriteString( t.Tiploc ).
    WriteString( t.Crs ).
    WriteString( t.Toc ).
    WriteString( t.Name ).
    WriteTime( t.Date )
}

func (t *Location) Read( c *codec.BinaryCodec ) {
  c.ReadString( &t.Tiploc ).
    ReadString( &t.Crs ).
    ReadString( &t.Toc ).
    ReadString( &t.Name ).
    ReadTime( &t.Date )
}
