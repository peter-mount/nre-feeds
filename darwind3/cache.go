package darwind3

import (
	bolt "github.com/etcd-io/bbolt"
	"log"
	"time"
)

const (
	alarmBucket       = "alarms"
	AssociationBucket = "assoc"
	messageBucket     = "messages"
	MetaBucket        = "meta"
	ScheduleBucket    = "schedule"
	TsBucket          = "ts"
)

func (d *DarwinD3) DBStatus() {
	//DBStatus(d.cache.db, alarmBucket, AssociationBucket, messageBucket, MetaBucket, ScheduleBucket, TsBucket)
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

func (r *DarwinD3) GetTimeMeta(n string) time.Time {
	v, _ := r.Meta.Get(n)
	if v == nil {
		return time.Time{}
	}
	return v.Data().(time.Time)
}

func (r *DarwinD3) PutTimeMeta(n string, o time.Time) {
	r.Meta.Add(n, o)
}

// Retrieve a schedule by it's rid
func (d *DarwinD3) GetSchedule(rid string) *Schedule {

	sched := d.GetScheduleNoResolve(rid)

	if sched == nil {
		sched = d.resolveSchedule(rid)
		d.PutSchedule(sched)
	}

	return sched
}

// Retrieve a schedule by it's rid
func (d *DarwinD3) PutSchedule(s *Schedule) {
	if s != nil {
		d.Schedules.Add(s.RID, s)
	}
}

func (d *DarwinD3) GetScheduleNoResolve(rid string) *Schedule {
	v, _ := d.Schedules.Get(rid)
	if v != nil {
		return v.Data().(*Schedule)
	}
	return nil
}
