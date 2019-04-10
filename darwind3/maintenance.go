package darwind3

import (
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/util"
	"log"
	"time"
)

// Expire old schedules after 3 days
const scheduleExpiry = time.Hour * 24 * 3

// PurgeSchedules purges all old schedules from the database, freeing up disk space
func (d *DarwinD3) PurgeSchedules() {
	// Run through the ts bucket and delete any entry older than our expiry time
	PurgeSchedules(d.cache.db, scheduleExpiry, DeleteSchedule)
}

func PurgeSchedules(db *bbolt.DB, maxAge time.Duration, del func(tx *bbolt.Tx, rid []byte)) {
	_ = db.Update(func(tx *bbolt.Tx) error {
		limit := time.Now().Add(-maxAge)
		log.Println("Expiring schedules older than ", limit.Format(util.HumanDateTime))

		ts := tx.Bucket([]byte(TsBucket))

		var t time.Time
		ec := 0
		dc := 0
		_ = ts.ForEach(func(k, v []byte) error {
			ec++
			if t.UnmarshalBinary(v) == nil && t.Before(limit) {
				del(tx, k)
				dc++
			}
			return nil
		})

		log.Printf("Expired %d/%d schedules\n", dc, ec)
		return nil
	})
}

// PurgeOrphans removes any schedules or ts entries which do not have a corresponding entry in the
// other bucket.
func (d *DarwinD3) PurgeOrphans() {
	PurgeOrphans(d.cache.db, DeleteSchedule)
}

func PurgeOrphans(db *bbolt.DB, del func(tx *bbolt.Tx, rid []byte)) {
	_ = db.Update(func(tx *bbolt.Tx) error {
		log.Println("Checking for orphans")

		ab := tx.Bucket([]byte(AssociationBucket))
		sb := tx.Bucket([]byte(ScheduleBucket))
		ts := tx.Bucket([]byte(TsBucket))

		c := testForOrphan(tx, del, ts, sb, ab)
		c += testForOrphan(tx, del, ab, ts, sb)
		c += testForOrphan(tx, del, sb, ab, ts)

		log.Printf("Removed %d orphans", c)
		return nil
	})
}

func testForOrphan(tx *bbolt.Tx, del func(tx *bbolt.Tx, rid []byte), targetBucket *bbolt.Bucket, srcBuckets ...*bbolt.Bucket) int {
	if targetBucket == nil {
		return 0
	}

	c := 0

	_ = targetBucket.ForEach(func(k, v []byte) error {
		found := true
		for _, bucket := range srcBuckets {
			if found && bucket != nil {
				found = bucket.Get(k) != nil
			}
		}
		if !found {
			del(tx, k)
			c++
		}
		return nil
	})

	return c
}
