package service

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/cron"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/nre-feeds/bin"
  "github.com/peter-mount/nre-feeds/darwinref"
  "github.com/peter-mount/nre-feeds/darwinupdate"
  "log"
)

type DarwinRefService struct {
  reference     darwinref.DarwinReference

  config       *bin.Config
  cron         *cron.CronService
  restService  *rest.Server
  updater      *darwinupdate.DarwinUpdate
}

func (a *DarwinRefService) Name() string {
  return "DarwinRefService"
}

func (a *DarwinRefService) Init( k *kernel.Kernel ) error {
  service, err := k.AddService( &bin.Config{} )
  if err != nil {
    return err
  }
  a.config = (service).(*bin.Config)

  service, err = k.AddService( &cron.CronService{} )
  if err != nil {
    return err
  }
  a.cron = (service).(*cron.CronService)

  service, err = k.AddService( &darwinupdate.DarwinUpdate{} )
  if err != nil {
    return err
  }
  a.updater = (service).(*darwinupdate.DarwinUpdate)

  service, err = k.AddService( &rest.Server{} )
  if err != nil {
    return err
  }
  a.restService = (service).(*rest.Server)

  // ReferenceUpdate
  return nil
}

func (a *DarwinRefService) PostInit() error {
  a.config.DbPath( &a.config.Database.Reference, "dwref.db" )
  if err := a.reference.OpenDB( a.config.Database.Reference ); err != nil {
    return err
  }

  if a.config.Ftp.Enabled {
    // Scheduled updates
    if a.config.Ftp.Schedule != "" {
      a.cron.AddFunc( a.config.Ftp.Schedule, func () {
        if err := a.updater.ReferenceUpdate( &a.reference ); err != nil {
          log.Println( "Failed import:", err )
        }
      })
      log.Println( "Auto Update using:", a.config.Ftp.Schedule )
    }

    // Initial import required?
    if a.updater.ImportRequiredTimetable( &a.reference ) {
      a.updater.ReferenceUpdate( &a.reference )
    }
  }

  // Rest services

  a.restService.Handle( "/reason/cancelled", a.reference.AllReasonCancelHandler ).Methods( "GET" )
  a.restService.Handle( "/reason/cancelled/{id}", a.reference.ReasonCancelHandler ).Methods( "GET" )

  a.restService.Handle( "/reason/late", a.reference.AllReasonLateHandler ).Methods( "GET" )
  a.restService.Handle( "/reason/late/{id}", a.reference.ReasonLateHandler ).Methods( "GET" )

  a.restService.Handle( "/via/{at}/{dest}/{loc1}", a.reference.ViaHandler ).Methods( "GET" )
  a.restService.Handle( "/via/{at}/{dest}/{loc1}/{loc2}", a.reference.ViaHandler ).Methods( "GET" )
  a.restService.Handle( "/via", a.reference.ViaResolveHandler ).Methods( "POST" )

  // Reference retrieval methods
  a.restService.Handle( "/crs/{id}", a.reference.CrsHandler ).Methods( "GET" )
  a.restService.Handle( "/crs", a.reference.AllCrsHandler ).Methods( "GET" )

  a.restService.Handle( "/tiploc", a.reference.TiplocsHandler ).Methods( "POST" )
  a.restService.Handle( "/tiploc/{id}", a.reference.TiplocHandler ).Methods( "GET" )

  a.restService.Handle( "/toc", a.reference.AllTocsHandler ).Methods( "GET" )
  a.restService.Handle( "/toc/{id}", a.reference.TocHandler ).Methods( "GET" )

  a.restService.Handle( "/search/{term}", a.reference.SearchHandler ).Methods( "GET" )

  return nil
}

func (a *DarwinRefService) Stop() {
  a.reference.Close()
}
