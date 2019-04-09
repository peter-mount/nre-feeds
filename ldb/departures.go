// LDB - Live Departure Boards
package ldb

import (
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
	"log"
	"time"
)

const (
	crsBucket      = "crs"
	messageBucket  = "message"
	scheduleBucket = "schedule"
	serviceBucket  = "service"
	tiplocBucket   = "tiploc"
)

type LDB struct {
	Darwin       string
	Reference    string
	EventManager *darwind3.DarwinEventManager
	db           *bolt.DB
}

type Task struct {
	d *LDB
	e *darwind3.DarwinEvent
	f func(*Task) error
}

func (d *LDB) Init(dbFile string) error {
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
		for _, bucket := range []string{crsBucket, scheduleBucket, serviceBucket, tiplocBucket} {
			err := d.createBucket(tx, bucket)
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

	// Ensure we have our stations loaded on startup
	d.RefreshStations()
	d.RequestStationMessages()
	d.DBStatus()

	return nil
}

func (d *LDB) DBStatus() {
	log.Printf("%-10s %8s %5s", "Bucket", "Keys", "Depth")
	_ = d.View(func(tx *bolt.Tx) error {
		for _, bucket := range []string{crsBucket, tiplocBucket} {
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

func (d *LDB) createBucket(tx *bolt.Tx, n string) error {
	key := []byte(n)
	b := tx.Bucket(key)
	if b == nil {
		log.Println("Creating bucket", n)
		_, err := tx.CreateBucket(key)
		return err
	}
	return nil
}

func (d *LDB) Update(f func(tx *bolt.Tx) error) error {
	return d.db.Update(f)
}

func (d *LDB) View(f func(tx *bolt.Tx) error) error {
	return d.db.View(f)
}
