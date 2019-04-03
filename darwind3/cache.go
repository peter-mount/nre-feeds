package darwind3

import (
	"encoding/json"
	bolt "github.com/etcd-io/bbolt"
	"time"
)

const (
	expiryTime  time.Duration = 120 * time.Second
	persistSize int           = 100
)

// Memory cache of schedules with disk persistance
type cache struct {
	db *bolt.DB
	tx *bolt.Tx
}

func (c *cache) initCache(cacheDir string) error {
	db, err := bolt.Open(cacheDir+"/schedules.dat", 0666, &bolt.Options{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		err := c.createBucket(tx, "schedule")
		if err != nil {
			return err
		}
		err = c.createBucket(tx, "meta")
		return err
	})
	if err != nil {
		return err
	}

	c.db = db
	return nil
}

func (c *cache) createBucket(tx *bolt.Tx, n string) error {
	key := []byte(n)
	b := tx.Bucket(key)
	if b == nil {
		_, err := tx.CreateBucket(key)
		return err
	}
	return nil
}

func (r *DarwinD3) GetMeta(n string, o interface{}) error {
	return r.View(func(tx *bolt.Tx) error {
		return r.GetMetaTx(tx, n, o)
	})
}

func (r *DarwinD3) GetMetaTx(tx *bolt.Tx, n string, o interface{}) error {
	b := tx.Bucket([]byte("meta")).Get([]byte(n))
	if b == nil {
		return nil
	}
	return json.Unmarshal(b, o)
}

func (r *DarwinD3) PutMeta(n string, o interface{}) error {
	return r.Update(func(tx *bolt.Tx) error {
		return r.PutMetaTx(tx, n, o)
	})
}

func (r *DarwinD3) PutMetaTx(tx *bolt.Tx, n string, o interface{}) error {
	b, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return tx.Bucket([]byte("meta")).Put([]byte(n), b)
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

		err := f(tx)

		r.Timetable = oldTT
		r.cache.tx = oldTx

		return err
	}

	if r.cache.tx != nil {
		return wrapper(r.cache.tx)
	}
	return r.Update(wrapper)
}

// Retrieve a schedule by it's rid
func (d *DarwinD3) GetSchedule(rid string) *Schedule {
	var sched *Schedule

	_ = d.View(func(tx *bolt.Tx) error {
		sched = d.cache.getSchedule(tx, rid)
		return nil
	})

	if sched == nil {
		sched = d.resolveSchedule(rid)
		if sched != nil {
			_ = d.Update(func(tx *bolt.Tx) error {
				d.cache.putSchedule(tx, sched)
				return nil
			})
		}
	}

	return sched
}

func (d *DarwinD3) GetScheduleNoResolve(rid string) *Schedule {
	var sched *Schedule

	_ = d.View(func(tx *bolt.Tx) error {
		sched = d.cache.getSchedule(tx, rid)
		return nil
	})

	return sched
}

func (d *cache) getSchedule(tx *bolt.Tx, rid string) *Schedule {
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
			ret = d.cache.putSchedule(tx, sched)
			return nil
		})
	} else {
		ret = d.cache.putSchedule(d.cache.tx, sched)
	}
	return ret
}

func (d *cache) putSchedule(tx *bolt.Tx, sched *Schedule) bool {
	key := []byte(sched.RID)
	sb := tx.Bucket([]byte("schedule"))
	b := sb.Get(key)
	if b != nil {
		os := ScheduleFromBytes(b)
		if os != nil && os.RID == sched.RID && !sched.Date.After(os.Date) {
			return false
		}
	}

	b, _ = sched.Bytes()
	err := sb.Put(key, b)
	if err != nil {
		return false
	}
	return true
}

// Delete a schedule
func (d *DarwinD3) DeleteSchedule(rid string) {
	_ = d.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("schedule")).Delete([]byte(rid))
	})
}
