package darwinref

import (
  bolt "github.com/coreos/bbolt"
  "encoding/xml"
  "github.com/peter-mount/golib/rest"
)

type CrsResponse struct {
  XMLName     xml.Name  `json:"-" xml:"crs"`
  Crs         string    `json:"crs" xml:"crs,attr"`
  Tiploc   []*Location  `json:"locations,omitempty" xml:"LocationRef"`
  Self        string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}

func (dr *DarwinReference) CrsHandler( r *rest.Rest ) error {
  return dr.View( func( tx *bolt.Tx ) error {
    crs := r.Var( "id" )

    if locations, exists := dr.GetCrs( tx, crs ); exists {
      resp := &CrsResponse{}
      r.Status( 200 ).Value( resp )

      resp.Crs = crs
      resp.Self = r.Self( r.Context() + "/crs/" + crs )

      resp.Tiploc = locations

      for _, l := range locations {
        l.SetSelf( r )
      }

    } else {
      r.Status( 404 )
    }

    return nil
  })
}
