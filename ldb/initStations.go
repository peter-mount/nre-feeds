package ldb

import (
  bolt "github.com/coreos/bbolt"
  "darwinref"
  "github.com/peter-mount/golib/codec"
  "log"
)
// initStations ensures we have all public stations defined on startup.
// Not doing so incurs a performance hit when a train references it for the
// first time.
func (d *LDB) initStations() {
  if err := d.Stations.Update( func() error {
    log.Println( "LDB: Initialising stations")

    if err := d.Reference.View( func( tx *bolt.Tx ) error {
      crsBucket := tx.Bucket( []byte( "DarwinCrs" ) )
      tiplocBucket := tx.Bucket( []byte( "DarwinTiploc" ) )

      return crsBucket.ForEach( func( k, v []byte ) error {
        var tpls []string
        codec.NewBinaryCodecFrom( v ).ReadStringArray( &tpls )

        var t []*darwinref.Location
        for _, tpl := range tpls {
          if loc, exists := d.Reference.GetTiplocBucket( tiplocBucket, tpl ); exists {
            t = append( t, loc )
          }
        }

        d.createStation( t )
        return nil
      })
    } ); err != nil {
      return err
    }

    log.Println( "LDB:", len( d.Stations.crs ), "Stations initialised")
    return nil
  } ); err != nil {
    log.Println( "LDB: Station import failed", err )
  }
}
