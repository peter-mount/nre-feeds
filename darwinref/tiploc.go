package darwinref

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/codec"
  //"time"
)

// Return a *Location for a tiploc
func (r *DarwinReference) GetTiploc( tx *bolt.Tx, tpl string ) ( *Location, bool ) {
  loc, exists := r.GetTiplocBucket( tx.Bucket( []byte("DarwinTiploc") ), tpl )
  return loc, exists
}

// Internal method that uses the shared writable transaction
func (r *DarwinReference) getTiploc( tpl string ) ( *Location, bool ) {
  loc, exists := r.GetTiplocBucket( r.tiploc, tpl )
  return loc, exists
}

func (r *DarwinReference) GetTiplocBucket( bucket *bolt.Bucket, tpl string ) ( *Location, bool ) {
  var loc *Location = &Location{}

  b := bucket.Get( []byte( tpl ) )
  if b == nil {
    return nil, false
  }

  codec.NewBinaryCodecFrom( b ).Read( loc )

  if( loc.Tiploc == "" ) {
    return nil,false
  }

  return loc, true
}

func (r *DarwinReference) addTiploc( loc *Location ) error{
  // Update only if it does not exist or is different
  if old, exists := r.getTiploc( loc.Tiploc ); !exists || !loc.Equals( old ) {
    codec := codec.NewBinaryCodec()
    codec.Write( loc )
    if codec.Error() != nil {
      return codec.Error()
    }

    return r.tiploc.Put( []byte( loc.Tiploc ), codec.Bytes() )
  }

  return nil
}
