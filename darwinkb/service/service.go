package service

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/nre-feeds/bin"
  "github.com/peter-mount/nre-feeds/darwinkb"
)

type DarwinKBService struct {
  darwinkb    *darwinkb.DarwinKB

  config       *bin.Config

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

  service, err = k.AddService( &rest.Server{} )
  if err != nil {
    return err
  }
  a.restService = (service).(*rest.Server)

  // ReferenceUpdate
  return nil
}

func (a *DarwinKBService) Start() error {

  a.restService.Handle( "/station/{crs}", a.StationHandler ).Methods( "GET" )

  // Expose the static directory so we offer the raw xml & full json files
  a.restService.Static( "/static/", a.config.KB.DataDir + "static/" )

  return nil
}
