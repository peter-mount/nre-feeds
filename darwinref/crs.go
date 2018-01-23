package darwinref

import (
  bolt "github.com/coreos/bbolt"
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

type crsimport struct {
   crs  map[string]map[string]string
   ref *DarwinReference
}

func (r *DarwinReference) newCrsImport() *crsimport {
  return &crsimport{
    crs: make( map[string]map[string]string ),
    ref: r,
  }
}

func (c *crsimport) append( loc *Location ) {
  if loc.Crs != "" {
    if e, ok := c.crs[ loc.Crs ]; ok {
      e[ loc.Tiploc ] = loc.Tiploc
    } else {
      c.crs[ loc.Crs ] = make( map[string]string )
      c.crs[ loc.Crs ][ loc.Tiploc ] = loc.Tiploc
    }
  }
}

func (c *crsimport) write() (error, int) {
  count := 0
  for k,v := range c.crs {
    if err, updated := c.writeCrs( k, v ); err != nil {
      return err, count
    } else if updated {
      count ++
    }
  }

  return nil, count
}

func (c *crsimport) writeCrs( k string, v map[string]string ) (error,bool) {
  // Now look at existing entry & skip if we match
  b := c.ref.crs.Get( []byte( k ) )
  if b != nil {
    var ar []string
    codec.NewBinaryCodecFrom( b ).ReadStringArray( &ar )

    if crsCompare( v, ar ) {
      return nil, false
    }
  }

  // get loaded tiplocs
  var tiplocs []string
  for t,_ := range v {
    tiplocs = append( tiplocs, t )
  }

  codec := codec.NewBinaryCodec()
  codec.WriteStringArray( tiplocs )
  if codec.Error() != nil {
    return codec.Error(), false
  }

  if err := c.ref.crs.Put( []byte( k ), codec.Bytes() ); err != nil {
    return err, false
  }

  return nil, true
}

func crsCompare( m map[string]string, ar []string ) bool {
  if len( ar ) == 0 {
    return false
  }

  for _, t := range ar {
    if _, exists := m[t]; !exists {
      return false
    }
  }
  return true
}
