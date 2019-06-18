package update

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/logger"
  "github.com/peter-mount/nre-feeds/bin"
  "github.com/peter-mount/nre-feeds/darwind3"
  "github.com/peter-mount/nre-feeds/darwinref/service"
)

type ReferenceUpdateService struct {
  config *bin.Config
  ref    *service.DarwinRefService
  logger *logger.LoggerService
}

func (a *ReferenceUpdateService) Name() string {
  return "ReferenceUpdateService"
}

func (a *ReferenceUpdateService) Init(k *kernel.Kernel) error {
  svce, err := k.AddService(&bin.Config{})
  if err != nil {
    return err
  }
  a.config = (svce).(*bin.Config)

  svce, err = k.AddService(&service.DarwinRefService{})
  if err != nil {
    return err
  }
  a.ref = (svce).(*service.DarwinRefService)

  svce, err = k.AddService(&logger.LoggerService{})
  if err != nil {
    return err
  }
  a.logger = (svce).(*logger.LoggerService)

  return nil
}

func (a *ReferenceUpdateService) Start() error {
  // Only listen if S3 is enabled
  if a.config.S3.Enabled {
    a.config.RabbitMQ.ConnectionName = "darwin ref"
    err := a.config.RabbitMQ.Connect()
    if err != nil {
      return err
    }

    em := darwind3.NewDarwinEventManager(&a.config.RabbitMQ, a.config.D3.EventKeyPrefix)
    err = em.ListenToEvents(darwind3.Event_TimeTableUpdate, a.referenceUpdateListener)
    if err != nil {
      return err
    }

    go a.findUpdates()
  }

  return nil
}
