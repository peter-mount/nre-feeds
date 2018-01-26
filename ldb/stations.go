package ldb

import (
  bolt "github.com/coreos/bbolt"
  "darwinref"
  "sync"
)

// Manages all stations
type Stations struct {
  mutex    *sync.Mutex
  // The stations being managed, by crs
  crs       map[string]*Station
  // The stations by tiploc
  tiploc    map[string]*Station
}

func NewStations() *Stations {
  s := &Stations{}
  s.crs = make( map[string]*Station )
  s.tiploc = make( map[string]*Station )
  s.mutex = &sync.Mutex{}
  return s
}

// Perform an action on the Stations instance with an exclusive lock
func (s *Stations) Update( f func() error ) error {
  s.mutex.Lock()
  defer s.mutex.Unlock()
  return f()
}

// GetStationCrs returns the Station instance by CRS or nil if not found
// Unlike GetStationTiploc this will not create a station if it's not found
func (d *LDB) GetStationCrs( crs string ) *Station {
  var station *Station

  d.Stations.Update( func() error {
    station = d.Stations.crs[ crs ]
    return nil
  } )

  return station
}

// GetStationTiploc returns the Station instance by Tiploc or nil if not found.
// Note: If we don't have an entry then this will create one
func (d *LDB) GetStationTiploc( tiploc string ) *Station {
  var station *Station

  // Perform read
  d.Stations.Update( func() error {
    station = d.Stations.tiploc[ tiploc ]
    return nil
  } )

  if station != nil {
    return station
  }

  // Still none so expensive but lock
  d.Stations.Update( func() error {
    station = d.Stations.tiploc[ tiploc ]

    if station != nil {
      return nil
    }

    var locs []*darwinref.Location
    d.Reference.View( func( tx *bolt.Tx ) error {
      // Lookup the tiploc
      loc, _ := d.Reference.GetTiploc( tx, tiploc )

      // Not found then bail - shouldn't happen unless reference data is out of sync
      if loc == nil {
        return nil
      }

      if loc.Crs == "" {
        // If no crs then use the single tiploc to prevent us from looking up again
        locs = append( locs, loc )
      } else {
        // Lookup by crs to get all of them
        locs, _ = d.Reference.GetCrs( tx, loc.Crs )
      }

      return nil
    } )

    if len( locs ) == 0 {
      return nil
    }

    station = d.createStation( locs )

    return nil
  } )

  return station
}

// Creates a station keyed by the supplied locations
func (d *LDB) createStation( locs []*darwinref.Location ) *Station {

  if len( locs ) == 0 {
    return nil
  }

  s := &Station{}
  s.Locations = locs

  // Mark public if we have a CRS & it doesn't start with X or Z
  s.public = locs[0].Crs != "" && locs[0].Crs[0] != 'X' && locs[0].Crs[0] != 'Z'

  if s.public {
    // Only public entries have a crs
    d.Stations.crs[ locs[0].Crs ] = s

    // Only public entries are useable so only create a mutex as needed
    s.mutex = &sync.Mutex{}
  }

  for _, l := range locs {
    d.Stations.tiploc[ l.Tiploc ] = s
  }

  return s
}
