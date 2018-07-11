package darwinkb

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/bolt"
  "github.com/peter-mount/golib/kernel/cron"
  "github.com/peter-mount/nre-feeds/bin"
  "os"
)

type DarwinKB struct {
  db      map[string]*bolt.BoltService
  config *bin.Config
  cron   *cron.CronService

  token   KBToken
}

func (r *DarwinKB) Name() string {
  return "DarwinKB"
}

var (
  buckets = []string{ "companies", "incidents", "serviceIndicators", "stations" }
)

func (a *DarwinKB) Init( k *kernel.Kernel ) error {

  service, err := k.AddService( &bin.Config{} )
  if err != nil {
    return err
  }
  a.config = (service).(*bin.Config)

  // see PostInit but use a dummy filename here, one db per bucket
  a.db = make( map[string]*bolt.BoltService)
  for _, bucket := range buckets {
    service, err = k.AddService( &bolt.BoltService{ FileName: bucket } )
    if err != nil {
      return err
    }
    a.db[bucket] = (service).(*bolt.BoltService)
  }

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
  for k, v := range a.db {
    v.FileName = a.config.KB.DataDir + "dwkb_" + k + ".db"
  }
  return nil
}

func (a *DarwinKB) Start() error {

  // Ensure the buckets exist in each db
  for k, db := range a.db {
    err := db.Update( func( tx *bolt.Tx ) error {
      _, err := tx.CreateBucketIfNotExists( k )
      return err
    } )
    if err != nil {
      return err
    }
  }

  // Check for updates during the morning & on startup
  a.cron.AddFunc( "0 30 4-9 * * *", a.refreshStations )
  a.refreshStations()

  // Incidents are regular intervals but not during the early hours
  a.cron.AddFunc( "0 0 0-1,5-23 * * *", a.refreshIncidents )
  a.refreshIncidents()

  // Refresh companies during the morning
  a.cron.AddFunc( "0 35 4-9 * * *", a.refreshCompanies )
  a.refreshCompanies()

  a.refreshServiceIndicators()
  
  return nil
}
