package darwinref

import (
  bolt "github.com/coreos/bbolt"
  "encoding/xml"
  "github.com/peter-mount/golib/codec"
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

func (dr *DarwinReference) AllCrsHandler( r *rest.Rest ) error {
  var t []*Location

  if err := dr.View( func( tx *bolt.Tx ) error {
    crsBucket := tx.Bucket( []byte( "DarwinCrs" ) )
    tiplocBucket := tx.Bucket( []byte( "DarwinTiploc" ) )

    return crsBucket.ForEach( func( k, v []byte ) error {
      var tpls []string
      codec.NewBinaryCodecFrom( v ).ReadStringArray( &tpls )

      for _, tpl := range tpls {
        if loc, exists := dr.GetTiplocBucket( tiplocBucket, tpl ); exists {
          loc.SetSelf( r )
          t = append( t, loc )
        }
      }

      r.Status( 200 ).Value( t )

      return nil
    } )
  }); err != nil {
    r.Status( 500 ).Value( err )
  }

  r.Status( 200 ).Value( t )
  return nil
}
