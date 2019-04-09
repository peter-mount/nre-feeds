package ldb

import (
	"bytes"
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
	"github.com/peter-mount/nre-feeds/util"
	"sort"
	"time"
)

const (
	dwellTime = time.Minute * 10
)

// Key must be unique so to support circular routes
// so we use Crs:RID:Tiploc:Time
func (s *Station) key(sched *darwind3.Schedule, idx int) []byte {
	loc := sched.Locations[idx]
	loc.UpdateTime()
	b := []byte(s.Crs + ":" + sched.RID + ":" + loc.Tiploc + ":")
	b = append(b, loc.Times.Time.Bytes()...)
	return b
}

func getServiceRID(b []byte) string {
	p := 0
	l := len(b)
	// Skip crs
	for p < l && b[p] != ':' {
		p++
	}
	// Find end of RID
	p++
	s := p
	for p < l && b[p] != ':' {
		p++
	}
	return string(b[s:p])
}

func getServiceTime(b []byte) *util.WorkingTime {
	p := 0
	l := len(b)
	// Skip crs, rid & tiploc
	for i := 0; i < 3; i++ {
		for p < l && b[p] != ':' {
			p++
		}
		p++
	}
	return util.WorkingTimeFromBytes(b[p:])
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
					_ = bucket.Put(key, b)
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
		prefix := []byte(s.Crs + ":" + rid + ":")
		bucket := tx.Bucket([]byte(serviceBucket))
		c := bucket.Cursor()
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			_ = bucket.Delete(k)
		}

	}

	return updated
}

func (s *Station) getServices(tx *bbolt.Tx, from *util.WorkingTime, to *util.WorkingTime) []*Service {
	var services []*Service

	if s.Public {

		// As a service can call at a station more than once, scan all and remove
		// every instance of it.
		prefix := []byte(s.Crs + ":")

		bucket := tx.Bucket([]byte(serviceBucket))
		c := bucket.Cursor()
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			t := getServiceTime(k)
			if t != nil && t.Between(from, to) {
				service := ServiceFromBytes(v)
				if !service.Location.Forecast.Departed {
					services = append(services, service)
				}
			}
		}

	}

	// sort into time order
	sort.SliceStable(services, func(i, j int) bool {
		return services[i].Compare(services[j])
	})

	return services
}
