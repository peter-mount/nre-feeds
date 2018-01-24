package darwinref

import (
  "encoding/xml"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/rest"
  "time"
)
// Defines a location, i.e. Station or passing point
type Location struct {
  XMLName     xml.Name  `json:"-" xml:"LocationRef"`
  Tiploc      string    `json:"tpl" xml:"tpl,attr"`
  Crs         string    `json:"crs,omitempty" xml:"crs,attr,omitempty"`
  Toc         string    `json:"toc,omitempty" xml:"toc,attr,omitempty"`
  Name        string    `json:"locname" xml:"locname,attr"`
  // Date entry was inserted into the database
  Date        time.Time `json:"date" xml:"date,attr"`
  // URL to this entity
  Self        string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}

// SetSelf sets the Self field to match this request
func (t *Location) SetSelf( r *rest.Rest ) {
  t.Self = r.Self( r.Context() + "/tiploc/" + t.Tiploc )
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
