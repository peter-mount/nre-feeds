package ldb

import (
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

func (d *LDB) getStationTiploc(tx *bbolt.Tx, tiploc string) *Station {
	crs, exists := d.tiplocs[tiploc]
	if !exists {
		return nil
	}
	return getStationCrs(tx, crs)
	/*
		key := []byte(tiploc)

			// Try to resolve the crs
		crs := tx.Bucket([]byte(tiplocBucket)).Get(key)
		if crs != nil {
			b := tx.Bucket([]byte(crsBucket)).Get([]byte(crs))
			if b != nil {
				return StationFromBytes(b)
			}
		}
		return nil
	*/
}

// GetStationTiploc returns the Station instance by Tiploc or nil if not found.
// Note: If we don't have an entry then this will create one
func (d *LDB) GetStationTiploc(tiploc string) *Station {
	crs, exists := d.tiplocs[tiploc]
	if !exists {
		return nil
	}
	return d.GetStationCrs(crs)
	/*
		var station *Station

			// Try to resolve the crs
		_ = d.View(func(tx *bbolt.Tx) error {
			station = d.getStationTiploc(tx, tiploc)
			return nil
		})

			return station
	*/
}

// Creates a station keyed by the supplied locations
func (d *LDB) createStation(tx *bbolt.Tx, locations []*darwinref.Location) *Station {

	if len(locations) == 0 {
		return nil
	}

	// Mark Public if we have a CRS & it doesn't start with X or Z
	crs := locations[0].Crs
	public := crs != "" && crs[0] != 'X' && crs[0] != 'Z'
	if !public {
		return nil
	}

	//tb := tx.Bucket([]byte(tiplocBucket))

	s := getStationCrs(tx, crs)
	if s == nil {
		s = &Station{}
		s.Crs = crs
		s.Locations = locations
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
			delete(d.tiplocs, t)
			//_ = tb.Delete([]byte(t))
		}
	}

	s.Public = public

	b, _ := s.Bytes()
	_ = tx.Bucket([]byte(crsBucket)).Put([]byte(crs), b)

	// Ensure all our tiplocs point to this crs
	//cb := []byte(crs)
	for _, l := range s.Locations {
		d.tiplocs[l.Tiploc] = crs
		//tpl := []byte(l.Tiploc)
		//b = tb.Get(tpl)
		//if b == nil || bytes.Compare(cb, b) != 0 {
		//	_ = tb.Put(tpl, cb)
		//}
	}

	return s
}
