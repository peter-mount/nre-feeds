// Handle the daily reference XML file from Darwin
package darwinref

import (
	"encoding/json"
	"errors"
	bolt "github.com/etcd-io/bbolt"
	"log"
	"time"
)

// Processed reference format
type DarwinReference struct {
	db *bolt.DB
	// Allow CIF.Close() to close the database.
	allowClose bool
	// Transaction used during import only
	tx *bolt.Tx
	// timetableId of the latest import
	timetableId string
	importDate  time.Time
	// Map of all locations by tiploc
	tiploc *bolt.Bucket
	// Map of all locations by CRS/3Alpha code
	crs *bolt.Bucket
	// Map of Toc (Operator) codes
	toc *bolt.Bucket
	// Reasons for a train being late
	lateRunningReasons *bolt.Bucket
	// Reasons for a train being cancelled at a location
	cancellationReasons *bolt.Bucket
	// CIS source
	cisSource map[string]string
	// via texts, map of at+","+ dest then array of possibilities
	via *bolt.Bucket
	// For speed, via's in memory
	viaMap *ViaMap
}

// OpenDB opens a DarwinReference database.
func (r *DarwinReference) OpenDB(dbFile string) error {
	if r.db != nil {
		return errors.New("DarwinReference Already attached to a Database")
	}

	if db, err := bolt.Open(dbFile, 0666, &bolt.Options{
		Timeout: 5 * time.Second,
	}); err != nil {
		return err
	} else {
		r.allowClose = true
		return r.useDB(db)
	}
}

// UseDB Allows an already open database to be used with DarwinReference.
func (r *DarwinReference) UseDB(db *bolt.DB) error {
	if r.db != nil {
		return errors.New("DarwinReference Already attached to a Database")
	}

	r.allowClose = false
	return r.useDB(db)
}

// common to OpenDB() && UseDB()
func (r *DarwinReference) useDB(db *bolt.DB) error {
	r.db = db

	// Now ensure the DB is initialised with the required buckets
	if err := r.initDB(); err != nil {
		return err
	}

	r.viaMap = r.NewViaMap()

	// Read metadata
	return r.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Meta")).Get([]byte("DarwinReference"))
		if b != nil {
			err := json.Unmarshal(b, r)
			if err != nil {
				return err
			}
		}

		if r.timetableId == "" {
			log.Println("DarwinReference needs importing")
		} else {
			log.Println("DarwinReference", r.timetableId, "imported", r.importDate)
		}

		return nil
	})
}

// Close the database.
// If OpenDB() was used to open the db then that db is closed.
// If UseDB() was used this simply detaches the DarwinReference from that DB. The DB is not closed()
func (r *DarwinReference) Close() {

	// Only close if we own the DB, e.g. via OpenDB()
	if r.allowClose && r.db != nil {
		r.db.Close()
	}

	// Detach
	r.db = nil
}

// Ensures we have the appropriate buckets
func (r *DarwinReference) initDB() error {

	buckets := []string{
		"Meta",
		"DarwinTiploc",
		"DarwinCrs",
		"DarwinToc",
		"DarwinLateReason",
		"DarwinCancelReason",
		"DarwinCIS",
		"DarwinVia"}

	return r.db.Update(func(tx *bolt.Tx) error {

		for _, n := range buckets {
			var nb []byte = []byte(n)
			if bucket := tx.Bucket(nb); bucket == nil {
				if _, err := tx.CreateBucket(nb); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// View performs a readonly operation on the database
func (r *DarwinReference) View(f func(*bolt.Tx) error) error {
	return r.db.View(f)
}

// Update performs a read write opertation on the database
func (r *DarwinReference) Update(f func(*bolt.Tx) error) error {
	return r.db.Update(f)
}

// internalUpdate is like Update but also sets our internal bucket references.
// This is usually used for importing data
func (r *DarwinReference) internalUpdate(f func(*bolt.Tx) error) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		r.tiploc = tx.Bucket([]byte("DarwinTiploc"))
		r.crs = tx.Bucket([]byte("DarwinCrs"))
		r.toc = tx.Bucket([]byte("DarwinToc"))
		r.cancellationReasons = tx.Bucket([]byte("DarwinCancelReason"))
		r.lateRunningReasons = tx.Bucket([]byte("DarwinLateReason"))
		r.via = tx.Bucket([]byte("DarwinVia"))

		return f(tx)
	})
}

// Return's the timetableId for this reference dataset
func (r *DarwinReference) TimetableId() string {
	return r.timetableId
}
