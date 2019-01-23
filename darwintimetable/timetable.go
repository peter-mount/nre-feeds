// Reference timetable
package darwintimetable

import (
  bolt "github.com/etcd-io/bbolt"
  "encoding/json"
  "errors"
  "log"
  "time"
)

type DarwinTimetable struct {
  db                   *bolt.DB
  // Allow CIF.Close() to close the database.
  allowClose            bool
  // Transaction used during import only
  tx                   *bolt.Tx
  // timetableId of the latest import
  timetableId           string              `xml:"timetableID,attr"`
  importDate            time.Time
  //Journeys      []*Journey            `xml:"Journey"`
  journeys        *bolt.Bucket
}

// OpenDB opens a DarwinReference database.
func (r *DarwinTimetable) OpenDB( dbFile string ) error {
  if r.db != nil {
    return errors.New( "Timetable Already attached to a Database" )
  }

  if db, err := bolt.Open( dbFile, 0666, &bolt.Options{
    Timeout: 5 * time.Second,
    } ); err != nil {
      return err
  } else {
    r.allowClose = true
    return r.useDB( db )
  }
}

// UseDB Allows an already open database to be used with DarwinReference.
func (r *DarwinTimetable) UseDB( db *bolt.DB ) error {
  if r.db != nil {
    return errors.New( "Timetable Already attached to a Database" )
  }

  r.allowClose = false
  return r.useDB( db )
}

// common to OpenDB() && UseDB()
func (r *DarwinTimetable) useDB( db *bolt.DB ) error {
  r.db = db

  // Now ensure the DB is initialised with the required buckets
  if err := r.initDB(); err != nil {
    return err
  }

  // Read metadata
  return r.View( func( tx *bolt.Tx ) error {
    b := tx.Bucket( []byte( "Meta" ) ).Get( []byte( "DarwinTimetable" ) )
    if b != nil {
      err := json.Unmarshal( b, r )
      if err != nil {
        return err
      }
    }

    if r.timetableId == "" {
      log.Println( "DarwinTimetable needs importing" )
    } else {
      log.Println( "DarwinTimetable", r.timetableId, "imported", r.importDate )
    }

    return nil
  })
}

// Close the database.
// If OpenDB() was used to open the db then that db is closed.
// If UseDB() was used this simply detaches the DarwinReference from that DB. The DB is not closed()
func (r *DarwinTimetable) Close() {

  // Only close if we own the DB, e.g. via OpenDB()
  if r.allowClose && r.db != nil {
    r.db.Close()
  }

  // Detach
  r.db = nil
}

// Ensures we have the appropriate buckets
func (r *DarwinTimetable) initDB() error {

  buckets := []string {
    "Meta",
    "DarwinAssoc",
    "DarwinJourney" }

  return r.db.Update( func( tx *bolt.Tx ) error {

    for _, n := range buckets {
      var nb []byte = []byte(n)
      if bucket := tx.Bucket( nb ); bucket == nil {
        if _, err := tx.CreateBucket( nb ); err != nil {
          return err
        }
      }
    }

    return nil
  })
}

// View performs a readonly operation on the database
func (r *DarwinTimetable) View( f func(*bolt.Tx) error ) error {
  return r.db.View( f )
}

// Update performs a read write opertation on the database
func (r *DarwinTimetable) Update( f func(*bolt.Tx) error ) error {
  return r.db.Update( f )
}

// internalUpdate is like Update but also sets our internal bucket references.
// This is usually used for importing data
func (r *DarwinTimetable) internalUpdate( f func(*bolt.Tx) error ) error {
  return r.db.Update( func( tx *bolt.Tx ) error {
    r.journeys = tx.Bucket( []byte("DarwinJourney") )

    return f(tx)
  })
}

// Return's the timetableId for this reference dataset
func (r *DarwinTimetable) TimetableId() string {
  return r.timetableId
}
