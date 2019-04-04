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
	_ = d.Update(func(tx *bbolt.Tx) error {

		limit := time.Now().Add(-scheduleExpiry)
		log.Println("Expiring schedules older than ", limit.Format(util.HumanDateTime))

		sb := tx.Bucket([]byte(scheduleBucket))
		ts := tx.Bucket([]byte(tsBucket))

		var t time.Time
		ec := 0
		dc := 0
		_ = ts.ForEach(func(k, v []byte) error {
			ec++
			if t.UnmarshalBinary(v) == nil && t.Before(limit) {
				_ = sb.Delete(k)
				_ = ts.Delete(k)
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
	_ = d.Update(func(tx *bbolt.Tx) error {

		log.Println("Checking for orphans")

		sb := tx.Bucket([]byte(scheduleBucket))
		ts := tx.Bucket([]byte(tsBucket))

		tc := 0
		sc := 0

		_ = ts.ForEach(func(k, v []byte) error {
			if sb.Get(k) == nil {
				_ = ts.Delete(k)
				tc++
			}
			return nil
		})

		_ = sb.ForEach(func(k, v []byte) error {
			if ts.Get(k) == nil {
				_ = sb.Delete(k)
				sc++
			}
			return nil
		})

		log.Printf("Removed %d ts & %d schedule orphans", tc, sc)
		return nil
	})
}
