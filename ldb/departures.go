// LDB - Live Departure Boards
package ldb

import (
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
	"time"
)

// The various bucket names
// crs holds the station details
// tiploc cross references tiploc to crs
// schedule holds the latest schedules
// message holds the station messages
// Also we use darwind3.ScheduleBucket & darwind3.TsBucket so we share the same mechanism's to keep
// schedules
const (
	crsBucket = "crs"
	//messageBucket = "message"
	serviceBucket = "service"
	//tiplocBucket  = "tiploc"
)

type LDB struct {
	Darwin       string
	Reference    string
	EventManager *darwind3.DarwinEventManager
	db           *bolt.DB
	tiplocs      map[string]string
	stations     map[string]*Station
}

type Task struct {
	d *LDB
	e *darwind3.DarwinEvent
	f func(*Task) error
}

func (d *LDB) Init(dbFile string) error {
	d.tiplocs = make(map[string]string)
	d.stations = make(map[string]*Station)

	db, err := bolt.Open(dbFile, 0666, &bolt.Options{
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
		for _, bucket := range []string{crsBucket, darwind3.ScheduleBucket, serviceBucket /*tiplocBucket,*/, darwind3.TsBucket} {
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
	d.db = db

	// Add listeners
	d.EventManager.ListenToEvents(darwind3.Event_ScheduleUpdated, d.locationListener)
	d.EventManager.ListenToEvents(darwind3.Event_Deactivated, d.deactivationListener)
	d.EventManager.ListenToEvents(darwind3.Event_StationMessage, d.stationMessageListener)

	// Ensure we have our stations loaded on startup, current messages & run the maintenance tasks
	d.RefreshStations()
	d.RequestStationMessages()
	d.PurgeSchedules()
	d.PurgeOrphans()
	d.PurgeServices()
	d.DBStatus()

	return nil
}

func (d *LDB) DBStatus() {
	darwind3.DBStatus(d.db, crsBucket, darwind3.ScheduleBucket, serviceBucket /*tiplocBucket,*/, darwind3.TsBucket)
}

func (d *LDB) Update(f func(tx *bolt.Tx) error) error {
	return d.db.Update(f)
}

func (d *LDB) View(f func(tx *bolt.Tx) error) error {
	return d.db.View(f)
}

func GetSchedule(tx *bolt.Tx, rid string) *darwind3.Schedule {
	sb := tx.Bucket([]byte("schedule"))
	b := sb.Get([]byte(rid))
	if b == nil {
		return nil
	}

	sched := darwind3.ScheduleFromBytes(b)
	if sched == nil || sched.RID == "" {
		return nil
	}

	return sched
}

func PutSchedule(tx *bolt.Tx, sched *darwind3.Schedule) bool {
	key := []byte(sched.RID)

	sb := tx.Bucket([]byte(darwind3.ScheduleBucket))

	b, _ := sched.Bytes()
	_ = sb.Put(key, b)

	b, err := sched.Date.MarshalBinary()
	if err == nil {
		_ = tx.Bucket([]byte(darwind3.TsBucket)).Put(key, b)
	}

	return true
}
