// Reference timetable
package darwintimetable

import (
  //bolt "github.com/etcd-io/bbolt"
  "encoding/xml"
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

type AssocService struct {
  RID       string    `json:"rid" xml:"rid,attr"`
  Wta       string    `json:"wta,omitempty" xml:"wta,attr"`
  Wtd       string    `json:"wtd,omitempty" xml:"wtd,attr"`
  Wtp       string    `json:"wtp,omitempty" xml:"wtp,attr"`
  Pta       string    `json:"pta,omitempty" xml:"pta,attr"`
  Ptd       string    `json:"ptd,omitempty" xml:"ptd,attr"`
}
