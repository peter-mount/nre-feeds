package ldb

import (
	"bytes"
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwinref"
)

func (d *LDB) PutStation(s *Station) {
	if s.Public {
		_ = d.Update(func(tx *bbolt.Tx) error {
			putStation(tx, s)
			return nil
		})
	}
}

func putStation(tx *bbolt.Tx, s *Station) {
	if s.Public {
		b, _ := s.Bytes()
		_ = tx.Bucket([]byte(crsBucket)).Put([]byte(s.Crs), b)
	}
}

// GetStationCrs returns the Station instance by CRS or nil if not found
// Unlike GetStationTiploc this will not create a station if it's not found
func (d *LDB) GetStationCrs(crs string) *Station {
	var station *Station

	_ = d.View(func(tx *bbolt.Tx) error {
		station = getStationCrs(tx, crs)
		return nil
	})

	return station
}

func getStationCrs(tx *bbolt.Tx, crs string) *Station {
	b := tx.Bucket([]byte(crsBucket)).Get([]byte(crs))
	if b != nil {
		return StationFromBytes(b)
	}
	return nil
}

// GetStationTiploc returns the Station instance by Tiploc or nil if not found.
// Note: If we don't have an entry then this will create one
func (d *LDB) GetStationTiploc(tiploc string) *Station {
	key := []byte(tiploc)

	var station *Station

	// Try to resolve the crs
	_ = d.View(func(tx *bbolt.Tx) error {
		crs := tx.Bucket([]byte(tiplocBucket)).Get(key)
		if crs != nil {
			b := tx.Bucket([]byte(crsBucket)).Get([]byte(crs))
			if b != nil {
				station = StationFromBytes(b)
			}
		}
		return nil
	})

	// Still none so expensive but lock
	/*
		d.Stations.Update(func() error {
			station = d.Stations.tiploc[tiploc]

			if station != nil {
				return nil
			}

			var locs []*darwinref.Location
			/ *
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
			* /

			if len(locs) == 0 {
				return nil
			}

			station = d.createStation(locs)

			return nil
		})
	*/

	return station
}

// Creates a station keyed by the supplied locations
func createStation(tx *bbolt.Tx, locations []*darwinref.Location) *Station {

	if len(locations) == 0 {
		return nil
	}

	// Mark Public if we have a CRS & it doesn't start with X or Z
	crs := locations[0].Crs
	public := crs != "" && crs[0] != 'X' && crs[0] != 'Z'
	if !public {
		return nil
	}

	tb := tx.Bucket([]byte(tiplocBucket))

	s := getStationCrs(tx, crs)
	if s == nil {
		s = &Station{}
		s.Crs = crs
		s.Locations = locations
		s.Services = make(map[string]*Service)
	} else {
		// Remove any tiplocs that have been removed
		tpl := make(map[string]interface{})
		for _, loc := range locations {
			tpl[loc.Tiploc] = true
		}
		for _, loc := range s.Locations {
			if _, exists := tpl[loc.Tiploc]; exists {
				delete(tpl, loc.Tiploc)
			}
		}
		for t, _ := range tpl {
			_ = tb.Delete([]byte(t))
		}
	}

	s.Public = public

	b, _ := s.Bytes()
	_ = tx.Bucket([]byte(crsBucket)).Put([]byte(crs), b)

	// Only Public entries are usable
	s.init()

	// Ensure all our tiplocs point to this crs
	cb := []byte(crs)
	for _, l := range s.Locations {
		tpl := []byte(l.Tiploc)
		b = tb.Get(tpl)
		if b == nil || bytes.Compare(cb, b) == 0 {
			_ = tb.Put(tpl, cb)
		}
	}

	return s
}
