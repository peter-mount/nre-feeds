package ldb

import (
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
)

func (d *LDB) PutSchedule(s *darwind3.Schedule) error {
	return d.Update(func(tx *bolt.Tx) error {
		darwind3.PutSchedule(tx, s)
		return nil
	})
}

func (d *LDB) GetSchedule(rid string) *darwind3.Schedule {
	var sched *darwind3.Schedule
	_ = d.View(func(tx *bolt.Tx) error {
		sched = darwind3.GetSchedule(tx, rid)
		return nil
	})
	return sched
}

func (d *LDB) RemoveSchedule(rid string) {
	_ = d.Update(func(tx *bolt.Tx) error {
		return removeSchedule(tx, rid)
	})
}

func removeSchedule(tx *bolt.Tx, rid string) error {
	darwind3.DeleteSchedule(tx, rid)

	bucket := tx.Bucket([]byte(serviceBucket))
	return bucket.ForEach(func(k, v []byte) error {
		if getServiceRID(k) == rid {
			return bucket.Delete(k)
		}
		return nil
	})
}
