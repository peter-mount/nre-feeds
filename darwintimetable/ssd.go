package darwintimetable

import (
  "encoding/json"
  "encoding/xml"
  "github.com/peter-mount/golib/codec"
  "time"
)

type SSD struct {
  t time.Time
}

func (a *SSD) Equals( b *SSD ) bool {
  return b != nil && a.t == b.t
}

func (t *SSD) Parse( s string ) {
  t.t, _ = time.Parse( "2006-01-02", s )
}

// Before is an SSD before a specified time
func (s *SSD) Before( t time.Time ) bool {
  return s.t.Before( t )
}

// BinaryCodec writer
func (t *SSD) Write( c *codec.BinaryCodec ) {
  c.WriteTime( t.t )
}

// BinaryCodec reader
func (t *SSD) Read( c *codec.BinaryCodec ) {
  c.ReadTime( &t.t )
}

// Custom JSON Marshaler.
func (t *SSD) MarshalJSON() ( []byte, error ) {
  return json.Marshal( t.String() )
}

// Custom XML Marshaler.
func (t *SSD) MarshalXMLAttr( name xml.Name ) ( xml.Attr, error ) {
  return xml.Attr{ Name: name, Value: t.String() }, nil
}

// String returns a SSD in "YYYY-MM-DD" format
func (t *SSD) String() string {
  return t.t.Format( "2006-01-02" )
}
