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

		sb := tx.Bucket([]byte(ScheduleBucket))
		ts := tx.Bucket([]byte(TsBucket))

		tc := 0
		sc := 0

		// Any TS entry with no schedule then delete it
		_ = ts.ForEach(func(k, v []byte) error {
			if sb.Get(k) == nil {
				del(tx, k)
				tc++
			}
			return nil
		})

		// Any Schedule entry with no ts then delete it
		_ = sb.ForEach(func(k, v []byte) error {
			if ts.Get(k) == nil {
				del(tx, k)
				sc++
			}
			return nil
		})

		log.Printf("Removed %d ts & %d schedule orphans", tc, sc)
		return nil
	})
}
