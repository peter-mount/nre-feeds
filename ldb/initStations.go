package ldb

import (
//  bolt "github.com/coreos/bbolt"
  "darwinref"
//  "github.com/peter-mount/golib/codec"
  "log"
)

// initStations ensures we have all public stations defined on startup.
// Not doing so incurs a performance hit when a train references it for the
// first time.
func (d *LDB) initStations() {
  if err := d.Stations.Update( func() error {
    log.Println( "LDB: Initialising stations")

    refClient := &darwinref.DarwinRefClient{ Url: d.Reference }

    if locations, err := refClient.GetStations(); err != nil {
      return err
    } else if locations != nil {
      // Map tiplocs by crs
      m := make( map[string][]*darwinref.Location )
      for _, loc := range locations {
        if s, ok := m[ loc.Crs ]; ok {
          m[ loc.Crs ] = append( s, loc )
        } else {
          s = make( []*darwinref.Location, 0 )
          m[ loc.Crs ] = append( s, loc )
        }
      }

      // create a station for each
      for _, l := range m {
        d.createStation( l )
      }
    }

    log.Println( "LDB:", len( d.Stations.crs ), "Stations initialised")
    return nil
  } ); err != nil {
    log.Println( "LDB: Station import failed", err )
  }

  //d.Darwin.ExpireStationMessages()
  //d.Darwin.BroadcastStationMessages()
}
