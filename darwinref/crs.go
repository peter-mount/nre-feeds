package darwinref

import (
  bolt "github.com/etcd-io/bbolt"
  "github.com/peter-mount/golib/codec"
)

// Return a *Location for a tiploc
func (r *DarwinReference) GetCrs( tx *bolt.Tx, t string ) ( []*Location, bool ) {
  loc, exists := r.GetCrsBucket( tx.Bucket( []byte("DarwinCrs") ), tx.Bucket( []byte("DarwinTiploc") ), t )
  return loc, exists
}

// Internal method that uses the shared writable transaction
func (r *DarwinReference) getCrs( t string ) ( []*Location, bool ) {
  loc, exists := r.GetCrsBucket( r.crs, r.tiploc, t )
  return loc, exists
}

func (r *DarwinReference) GetCrsBucket( crsbucket *bolt.Bucket, tiplocbucket *bolt.Bucket, crs string ) ( []*Location, bool ) {
  b := crsbucket.Get( []byte( crs ) )
  if b == nil {
    return nil, false
  }

  var ar []string
  codec.NewBinaryCodecFrom( b ).ReadStringArray( &ar )

  if len( ar ) == 0 {
    return nil, false
  }

  var t []*Location
  for _, k := range ar {
    if loc, exists := r.GetTiplocBucket( tiplocbucket, k ); exists {
      t = append( t, loc )
    }
  }

  return t, len( t ) > 0
}
