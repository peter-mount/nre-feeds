package ldb

import (
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
)

func (d *LDB) PutSchedule(s *darwind3.Schedule) error {
	return d.Update(func(tx *bolt.Tx) error {
		return putSchedule(tx, s)
	})
}

func putSchedule(tx *bolt.Tx, s *darwind3.Schedule) error {
	b, err := s.Bytes()
	if err != nil {
		return err
	}
	return tx.Bucket([]byte(scheduleBucket)).Put([]byte(s.RID), b)
}

func (d *LDB) GetSchedule(rid string) *darwind3.Schedule {
	var sched *darwind3.Schedule
	_ = d.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(scheduleBucket)).Get([]byte(rid))
		if b != nil {
			sched = darwind3.ScheduleFromBytes(b)
		}
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
	err := tx.Bucket([]byte(scheduleBucket)).Delete([]byte(rid))
	if err != nil {
		return err
	}

	bucket := tx.Bucket([]byte(serviceBucket))
	return bucket.ForEach(func(k, v []byte) error {
		if getServiceRID(k) == rid {
			return bucket.Delete(k)
		}
		return nil
	})
}
