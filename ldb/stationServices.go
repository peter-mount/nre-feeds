package ldb

import (
	"bytes"
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
	"time"
)

const (
	dwellTime = time.Minute * 10
)

// Key must be unique so to support circular routes
// so we use the RID, tiploc and the timetable time
func (s *Station) key(sched *darwind3.Schedule, idx int) []byte {
	loc := sched.Locations[idx]
	return []byte(sched.RID + ":" + loc.Tiploc + ":" + loc.Times.Time.String())
}

// Adds a service to the station
func (s *Station) addService(tx *bbolt.Tx, e *darwind3.DarwinEvent, idx int) bool {
	// Only Public stations can be updated. Pass to the channel so the worker thread
	// can read it
	if s.Public && e.Schedule.Locations[idx].Times.IsPublic() {
		// Add service

		t := e.Schedule.GetTime(idx)
		if t.After(time.Now().Add(-dwellTime)) {
			key := s.key(e.Schedule, idx)

			bucket := tx.Bucket([]byte(serviceBucket))

			service := ServiceFromBytes(bucket.Get(key))
			if service == nil || service.Date.Before(t) {
				if service == nil {
					service = &Service{}
				}

				if service.update(e.Schedule, idx) {
					b, _ := service.Bytes()
					bucket.Put(key, b)
					return true
				}
			}
		}

	}
	return false
}

func (s *Station) removeDepartedService(tx *bbolt.Tx, e *darwind3.DarwinEvent, idx int) bool {
	b := tx.Bucket([]byte(serviceBucket)).Get(s.key(e.Schedule, idx))

	return b != nil
}

// Removes all entries of a service from a station.
func (s *Station) removeService(tx *bbolt.Tx, rid string) bool {
	updated := false

	if s.Public {

		// As a service can call at a station more than once, scan all and remove
		// every instance of it.
		prefix := []byte(rid + ":")
		bucket := tx.Bucket([]byte(serviceBucket))
		c := bucket.Cursor()
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			bucket.Delete(k)
		}

	}

	return updated
}
