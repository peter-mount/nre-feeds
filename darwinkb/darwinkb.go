package darwinkb

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/bolt"
  "github.com/peter-mount/nre-feeds/bin"
  "os"
)

type DarwinKB struct {
  db           *bolt.BoltService

  config       *bin.Config
  boltDb       *bolt.BoltService

  token         KBToken
}

func (r *DarwinKB) Name() string {
  return "DarwinKB"
}

// OpenDB opens a DarwinReference database.
func (r *DarwinKB) OpenDB( db *bolt.BoltService ) {
  r.db = db
}

func (a *DarwinKB) Init( k *kernel.Kernel ) error {

  service, err := k.AddService( &bin.Config{} )
  if err != nil {
    return err
  }
  a.config = (service).(*bin.Config)

  // see PostInit but use a dummy filename here
  service, err = k.AddService( &bolt.BoltService{ FileName: "dwkb.db" } )
  if err != nil {
    return err
  }
  a.boltDb = (service).(*bolt.BoltService)

  return nil
}

func (a *DarwinKB) PostInit() error {
  if a.config.KB.DataDir == "" {
    a.config.KB.DataDir = "/database/"
  }
  err := os.MkdirAll( a.config.KB.DataDir + "static/", 0x755 )
  if err != nil {
    return err
  }

  // This will work as the db isn't stated yet
  a.boltDb.FileName = a.config.KB.DataDir + "kb.db"
  return nil
}

func (a *DarwinKB) Start() error {

  err := a.boltDb.Update( func( tx *bolt.Tx ) error {
    _, err := tx.CreateBucketIfNotExists( "stations" )
    return err
  } )
  if err != nil {
    return err
  }

  err = a.refreshStations()
  if err != nil {
    return err
  }


  return nil
}
