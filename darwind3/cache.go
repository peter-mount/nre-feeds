package darwind3

import (
	"encoding/json"
	bolt "github.com/etcd-io/bbolt"
	"log"
	"time"
)

// Memory cache of schedules with disk persistance
type cache struct {
	db *bolt.DB
	tx *bolt.Tx
}

const (
	alarmBucket       = "alarms"
	AssociationBucket = "assoc"
	messageBucket     = "messages"
	MetaBucket        = "meta"
	ScheduleBucket    = "schedule"
	TsBucket          = "ts"
)

func (c *cache) initCache(cacheDir string) error {
	db, err := bolt.Open(cacheDir, 0666, &bolt.Options{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}

	// Ensure our buckets exist
	// meta for metadata
	// schedule for the live data
	// ts for the times per rid - used for cleaning up
	err = db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range []string{alarmBucket, AssociationBucket, messageBucket, MetaBucket, ScheduleBucket, TsBucket} {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	c.db = db
	return nil
}

func (d *DarwinD3) DBStatus() {
	DBStatus(d.cache.db, alarmBucket, AssociationBucket, messageBucket, MetaBucket, ScheduleBucket, TsBucket)
}

func DBStatus(db *bolt.DB, buckets ...string) {
	log.Printf("%-10s %8s %5s", "Bucket", "Keys", "Depth")
	_ = db.View(func(tx *bolt.Tx) error {
		for _, bucket := range buckets {
			bs := tx.Bucket([]byte(bucket)).
				Stats()

			log.Printf(
				"%-10s %8d %5d",
				bucket,
				bs.KeyN,
				bs.Depth,
			)
		}
		return nil
	})
}

func (r *DarwinD3) GetMeta(n string, o interface{}) error {
	return r.View(func(tx *bolt.Tx) error {
		return GetMeta(tx, n, o)
	})
}

func GetMeta(tx *bolt.Tx, n string, o interface{}) error {
	b := tx.Bucket([]byte(MetaBucket)).Get([]byte(n))
	if b == nil {
		return nil
	}
	return json.Unmarshal(b, o)
}

func (r *DarwinD3) PutMeta(n string, o interface{}) error {
	return r.Update(func(tx *bolt.Tx) error {
		return PutMeta(tx, n, o)
	})
}

func PutMeta(tx *bolt.Tx, n string, o interface{}) error {
	b, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return tx.Bucket([]byte(MetaBucket)).Put([]byte(n), b)
}

// View performs a readonly operation on the database
func (r *DarwinD3) View(f func(*bolt.Tx) error) error {
	return r.cache.db.View(f)
}

// SnapshotUpdate performs a read write opertation on the database
func (r *DarwinD3) Update(f func(*bolt.Tx) error) error {
	return r.cache.db.Update(f)
}

func (r *DarwinD3) BulkUpdate(f func(*bolt.Tx) error) error {
	wrapper := func(tx *bolt.Tx) error {
		oldTx := r.cache.tx
		r.cache.tx = tx

		oldTT := r.Timetable
		r.Timetable = ""

		defer func() {
			r.Timetable = oldTT
			r.cache.tx = oldTx
		}()

		return f(tx)
	}

	if r.cache.tx != nil {
		return wrapper(r.cache.tx)
	}
	return r.Update(wrapper)
}

// Retrieve a schedule by it's rid
func (d *DarwinD3) GetSchedule(rid string) *Schedule {

	sched := d.GetScheduleNoResolve(rid)

	if sched == nil {
		sched = d.resolveSchedule(rid)
		if sched != nil {
			_ = d.Update(func(tx *bolt.Tx) error {
				PutSchedule(tx, sched)
				return nil
			})
		}
	}

	return sched
}

func (d *DarwinD3) GetScheduleNoResolve(rid string) *Schedule {
	var sched *Schedule

	_ = d.View(func(tx *bolt.Tx) error {
		sched = GetSchedule(tx, rid)
		return nil
	})

	return sched
}

func GetSchedule(tx *bolt.Tx, rid string) *Schedule {
	sb := tx.Bucket([]byte("schedule"))
	b := sb.Get([]byte(rid))
	if b == nil {
		return nil
	}

	sched := ScheduleFromBytes(b)
	if sched == nil || sched.RID == "" {
		return nil
	}

	return sched
}

// Store a schedule by it's rid
func (d *DarwinD3) PutSchedule(sched *Schedule) bool {
	ret := false

	if d.cache.tx == nil {
		_ = d.Update(func(tx *bolt.Tx) error {
			ret = PutSchedule(tx, sched)
			return nil
		})
	} else {
		ret = PutSchedule(d.cache.tx, sched)
	}
	return ret
}

func PutSchedule(tx *bolt.Tx, sched *Schedule) bool {
	key := []byte(sched.RID)

	sb := tx.Bucket([]byte(ScheduleBucket))

	/*
			TODO remove this if neccessary
		b := sb.Get(key)
		if b != nil {
			os := ScheduleFromBytes(b)
			if os != nil && os.RID == sched.RID && !sched.Date.After(os.Date) {
				return false
			}
		}
	*/

	b, _ := sched.Bytes()
	_ = sb.Put(key, b)

	b, err := sched.Date.MarshalBinary()
	if err == nil {
		_ = tx.Bucket([]byte(TsBucket)).Put(key, b)
	}

	return true
}

// Delete a schedule
func (d *DarwinD3) DeleteSchedule(rid string) {
	_ = d.Update(func(tx *bolt.Tx) error {
		DeleteSchedule(tx, []byte(rid))
		return nil
	})
}

func DeleteSchedule(tx *bolt.Tx, rid []byte) {
	deleteSchedule(tx.Bucket([]byte(AssociationBucket)), rid)
	deleteSchedule(tx.Bucket([]byte(ScheduleBucket)), rid)
	deleteSchedule(tx.Bucket([]byte(TsBucket)), rid)
}

func deleteSchedule(bucket *bolt.Bucket, rid []byte) {
	if bucket != nil {
		_ = bucket.Delete(rid)
	}
}
