package darwinkb

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/bolt"
  "github.com/peter-mount/golib/kernel/cron"
  "github.com/peter-mount/nre-feeds/bin"
  "os"
)

type DarwinKB struct {
  db           *bolt.BoltService

  boltDb       *bolt.BoltService
  config       *bin.Config
  cron         *cron.CronService

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

  service, err = k.AddService( &cron.CronService{} )
  if err != nil {
    return err
  }
  a.cron = (service).(*cron.CronService)

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
  a.boltDb.FileName = a.config.KB.DataDir + "dwkb.db"
  return nil
}

func (a *DarwinKB) Start() error {

  // Ensure the buckets exist
  err := a.boltDb.Update( func( tx *bolt.Tx ) error {
    buckets :=[]string{ "stations" }
    for _, n := range buckets {
      _, err := tx.CreateBucketIfNotExists( n )
      if err != nil {
        return err
      }
    }
    return nil
  } )
  if err != nil {
    return err
  }

  // Check for updates during the morning & on startup
  a.cron.AddFunc( "0 30 4-9 * * *", a.refreshStations )
  a.refreshStations()

  return nil
}
