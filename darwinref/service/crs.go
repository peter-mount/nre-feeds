package service

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/nre-feeds/darwinref"
)

func (dr *DarwinRefService) CrsHandler( r *rest.Rest ) error {
  return dr.reference.View( func( tx *bolt.Tx ) error {
    crs := r.Var( "id" )

    if locations, exists := dr.reference.GetCrs( tx, crs ); exists {
      resp := &darwinref.CrsResponse{}
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

func (dr *DarwinRefService) AllCrsHandler( r *rest.Rest ) error {
  var t []*darwinref.Location

  if err := dr.reference.View( func( tx *bolt.Tx ) error {
    crsBucket := tx.Bucket( []byte( "DarwinCrs" ) )
    tiplocBucket := tx.Bucket( []byte( "DarwinTiploc" ) )

    return crsBucket.ForEach( func( k, v []byte ) error {
      var tpls []string
      codec.NewBinaryCodecFrom( v ).ReadStringArray( &tpls )

      for _, tpl := range tpls {
        if loc, exists := dr.reference.GetTiplocBucket( tiplocBucket, tpl ); exists {
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