package update

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/nre-feeds/bin"
  "github.com/peter-mount/nre-feeds/darwind3"
  "github.com/peter-mount/nre-feeds/darwintimetable/service"
  "log"
)

type TimetableUpdateService struct {
  config       *bin.Config
  timetable    *service.DarwinTimetableService
}

func (a *TimetableUpdateService) Name() string {
  return "TimetableUpdateService"
}

func (a *TimetableUpdateService) Init( k *kernel.Kernel ) error {
  svce, err := k.AddService( &bin.Config{} )
  if err != nil {
    return err
  }
  a.config = (svce).(*bin.Config)

  svce, err = k.AddService( &service.DarwinTimetableService{} )
  if err != nil {
    return err
  }
  a.timetable = (svce).(*service.DarwinTimetableService)

  return nil
}

func (a *TimetableUpdateService) Start() error {
  a.config.RabbitMQ.ConnectionName = "darwin tt"
  a.config.RabbitMQ.Connect()

  em := darwind3.NewDarwinEventManager( &a.config.RabbitMQ, a.config.D3.EventKeyPrefix )
  em.ListenToEvents( darwind3.Event_TimeTableUpdate, a.timetableUpdateListener )

  // debug only
  log.Println( "TUS started" )

  return nil
}
