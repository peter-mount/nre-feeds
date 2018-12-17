// Reference timetable
package darwintimetable

import (
  //bolt "github.com/etcd-io/bbolt"
  "encoding/xml"
  "github.com/peter-mount/golib/codec"
  //"github.com/peter-mount/golib/rest"
  "time"
)

type Association struct {
  XMLName   xml.Name      `json:"-" xml:"Association"`
  Main      AssocService  `json:"main" xml:"main"`
  Assoc     AssocService  `json:"assoc" xml:"assoc"`
  Tiploc    string        `json:"tiploc" xml:"tiploc,attr"`
  Category  string        `json:"category" xml:"category,attr"`
  Cancelled bool          `json:"cancelled" xml:"isCancelled,attr"`
  Deleted   bool          `json:"deleted" xml:"isDeleted,attr"`
  // Date entry was inserted into the database
  Date        time.Time   `json:"date" xml:"date,attr"`
  // URL to this entity
  Self        string      `json:"self" xml:"self,attr,omitempty"`
}

func (a *Association) Equals( b *Association ) bool {
  if b == nil {
    return false
  }
  return a.Main.RID == b.Main.RID &&
    a.Assoc.RID == b.Assoc.RID &&
    a.Tiploc == b.Tiploc &&
    a.Category == b.Category &&
    a.Cancelled == b.Cancelled &&
    a.Deleted == b.Deleted
}

func (t *Association) Write( c *codec.BinaryCodec ) {
  c.Write( &t.Main ).
    Write( &t.Assoc ).
    WriteString( t.Tiploc ).
    WriteString( t.Category ).
    WriteBool( t.Cancelled ).
    WriteBool( t.Deleted ).
    WriteTime( t.Date )
}

func (t *Association) Read( c *codec.BinaryCodec ) {
  c.Read( &t.Main ).
    Read( &t.Assoc ).
    ReadString( &t.Tiploc ).
    ReadString( &t.Category ).
    ReadBool( &t.Cancelled ).
    ReadBool( &t.Deleted ).
    ReadTime( &t.Date )
}

type AssocService struct {
  RID       string    `json:"rid" xml:"rid,attr"`
  Wta       string    `json:"wta" xml:"wta,attr"`
  Wtd       string    `json:"wtd" xml:"wtd,attr"`
  Wtp       string    `json:"wtp" xml:"wtp,attr"`
  Pta       string    `json:"pta" xml:"pta,attr"`
  Ptd       string    `json:"ptd" xml:"ptd,attr"`
}

func (t *AssocService) Write( c *codec.BinaryCodec ) {
  c.WriteString( t.RID ).
    WriteString( t.Wta ).
    WriteString( t.Wtd ).
    WriteString( t.Wtp ).
    WriteString( t.Pta ).
    WriteString( t.Ptd )
}

func (t *AssocService) Read( c *codec.BinaryCodec ) {
  c.ReadString( &t.RID ).
    ReadString( &t.Wta ).
    ReadString( &t.Wtd ).
    ReadString( &t.Wtp ).
    ReadString( &t.Pta ).
    ReadString( &t.Ptd )
}
