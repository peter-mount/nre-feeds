package ldb

import (
	"bytes"
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
	"log"
	"time"
)

func (d *LDB) PutSchedule(s *darwind3.Schedule) error {
	return d.Update(func(tx *bolt.Tx) error {
		PutSchedule(tx, s)
		return nil
	})
}

func (d *LDB) GetSchedule(rid string) *darwind3.Schedule {
	var sched *darwind3.Schedule
	_ = d.View(func(tx *bolt.Tx) error {
		sched = GetSchedule(tx, rid)
		return nil
	})
	return sched
}

func (d *LDB) RemoveSchedule(rid string) {
	_ = d.Update(func(tx *bolt.Tx) error {
		RemoveSchedule(tx, []byte(rid))
		return nil
	})
}

func DeleteSchedule(tx *bolt.Tx, rid []byte) {
	deleteSchedule(tx.Bucket([]byte(darwind3.AssociationBucket)), rid)
	deleteSchedule(tx.Bucket([]byte(darwind3.ScheduleBucket)), rid)
	deleteSchedule(tx.Bucket([]byte(darwind3.TsBucket)), rid)
}

func deleteSchedule(bucket *bolt.Bucket, rid []byte) {
	if bucket != nil {
		_ = bucket.Delete(rid)
	}
}

func RemoveSchedule(tx *bolt.Tx, rid []byte) {
	/*darwind3.*/ DeleteSchedule(tx, rid)

	bucket := tx.Bucket([]byte(serviceBucket))
	_ = bucket.ForEach(func(k, v []byte) error {
		if bytes.Compare(getServiceRID(k), rid) == 0 {
			return bucket.Delete(k)
		}
		return nil
	})
}

func (d *LDB) PurgeSchedules() {
	darwind3.PurgeSchedules(d.db, 24*time.Hour, RemoveSchedule)
}

func (d *LDB) PurgeOrphans() {
	darwind3.PurgeOrphans(d.db, RemoveSchedule)
}

// PurgeServices looks for any services who's schedule has been deleted
func (d *LDB) PurgeServices() {
	_ = d.Update(func(tx *bolt.Tx) error {
		log.Println("Looking for orphaned services")

		schedBucket := tx.Bucket([]byte(darwind3.ScheduleBucket))
		svcBucket := tx.Bucket([]byte(serviceBucket))

		deleted := 0
		count := 0

		_ = svcBucket.ForEach(func(k, v []byte) error {
			rid := getServiceRID(k)
			count++
			if schedBucket.Get([]byte(rid)) == nil {
				_ = svcBucket.Delete(k)
				deleted++
			}
			return nil
		})

		log.Printf("Purged %d/%d services", deleted, count)
		return nil
	})
}
