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
  // Tiploc of this location
  Tiploc      string    `json:"tpl" xml:"tpl,attr"`
  // CRS of this station, "" for none
  Crs         string    `json:"crs,omitempty" xml:"crs,attr,omitempty"`
  // TOC who manages this station
  Toc         string    `json:"toc,omitempty" xml:"toc,attr,omitempty"`
  // Name of this station
  Name        string    `json:"locname" xml:"locname,attr"`
  // True if this represents a Station, Bus stop or Ferry Terminal
  // i.e. Crs is present but does not start with X or Z
  Station     bool      `json:"station" xml:"station,attr"`
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
  t.Station = t.Crs != "" && !(t.Crs[0]=='X'||t.Crs[0]=='Z')
}

// IsPublic returns true if this Location represents a public station.
// This is defined as having a Crs and one that does not start with X (Network Rail,
// some Bus stations and some Ferry terminals) and Z (usually London Underground)
func (t *Location) IsPublic() bool {
  return t.Crs != "" && t.Crs[0] != 'X' && t.Crs[0] != 'Z'
}
