package service

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/cron"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/nre-feeds/bin"
  "github.com/peter-mount/nre-feeds/darwinkb"
)

type DarwinKBService struct {
  darwinkb    *darwinkb.DarwinKB

  config       *bin.Config

  cron         *cron.CronService
  restService  *rest.Server
}

func (a *DarwinKBService) Name() string {
  return "DarwinKBService"
}

func (a *DarwinKBService) Init( k *kernel.Kernel ) error {

  service, err := k.AddService( &bin.Config{} )
  if err != nil {
    return err
  }
  a.config = (service).(*bin.Config)

  service, err = k.AddService( &darwinkb.DarwinKB{} )
  if err != nil {
    return err
  }
  a.darwinkb = (service).(*darwinkb.DarwinKB)

  service, err = k.AddService( &cron.CronService{} )
  if err != nil {
    return err
  }
  a.cron = (service).(*cron.CronService)

  service, err = k.AddService( &rest.Server{} )
  if err != nil {
    return err
  }
  a.restService = (service).(*rest.Server)

  // ReferenceUpdate
  return nil
}

func (a *DarwinKBService) Start() error {
  // Expire old messages every 15 minutes & run an expire on startup
  //a.cron.AddFunc( "0 0/15 * * * *", a.darwind3.ExpireStationMessages )
  //go a.darwind3.ExpireStationMessages()

  return nil
}
