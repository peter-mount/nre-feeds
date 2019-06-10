package darwind3

import (
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/golib/statistics"
	"github.com/peter-mount/nre-feeds/util"
	"log"
	"runtime"
	"runtime/debug"
	"time"
)

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

func GC() {
	log.Println("GC")
	gcStats()
	debug.FreeOSMemory()
	gcStats()
}

func gcStats() {
	var gcStats debug.GCStats

	debug.ReadGCStats(&gcStats)
	log.Printf(
		"GC: %d Last %s Paused %s\n",
		gcStats.NumGC,
		gcStats.LastGC.Format(util.HumanDateTime),
		gcStats.PauseTotal.String(),
	)
}

func SubmitMemStats(prefix string) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	statistics.Set(prefix+".mem.heap.nextgc", int64(memStats.NextGC))
	statistics.Set(prefix+".mem.heap.alloc", int64(memStats.HeapAlloc))
	statistics.Set(prefix+".mem.heap.sys", int64(memStats.HeapSys))
	statistics.Set(prefix+".mem.heap.idle", int64(memStats.HeapIdle))
	statistics.Set(prefix+".mem.heap.inuse", int64(memStats.HeapInuse))
	statistics.Set(prefix+".mem.heap.released", int64(memStats.HeapReleased))
	statistics.Set(prefix+".mem.heap.objects", int64(memStats.HeapObjects))

	statistics.Set(prefix+".mem.stack.inuse", int64(memStats.StackInuse))
	statistics.Set(prefix+".mem.stack.sys", int64(memStats.StackSys))
}
