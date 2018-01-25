// darwind3 handles the real time push port feed
package darwind3

import (
  bolt "github.com/coreos/bbolt"
  "errors"
  "time"
)

type DarwinD3 struct {
  //
  db                   *bolt.DB
  // Allow CIF.Close() to close the database.
  allowClose            bool
  // Transaction used during import only
  tx                   *bolt.Tx
}

// OpenDB opens a DarwinReference database.
func (r *DarwinD3) OpenDB( dbFile string ) error {
  if r.db != nil {
    return errors.New( "DarwinReference Already attached to a Database" )
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
func (r *DarwinD3) UseDB( db *bolt.DB ) error {
  if r.db != nil {
    return errors.New( "DarwinReference Already attached to a Database" )
  }

  r.allowClose = false
  return r.useDB( db )
}

// common to OpenDB() && UseDB()
func (r *DarwinD3) useDB( db *bolt.DB ) error {
  r.db = db

  // Now ensure the DB is initialised with the required buckets
  if err := r.initDB(); err != nil {
    return err
  }

  return nil
}

// Close the database.
// If OpenDB() was used to open the db then that db is closed.
// If UseDB() was used this simply detaches the DarwinReference from that DB. The DB is not closed()
func (r *DarwinD3) Close() {

  // Only close if we own the DB, e.g. via OpenDB()
  if r.allowClose && r.db != nil {
    r.db.Close()
  }

  // Detach
  r.db = nil
}

// Ensures we have the appropriate buckets
func (r *DarwinD3) initDB() error {

  buckets := []string {
    "Meta",
    "DarwinRID" }

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
func (r *DarwinD3) View( f func(*bolt.Tx) error ) error {
  return r.db.View( f )
}

// Update performs a read write opertation on the database
func (r *DarwinD3) Update( f func(*bolt.Tx) error ) error {
  return r.db.Update( f )
}
