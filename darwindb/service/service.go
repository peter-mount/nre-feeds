package service

import (
	"github.com/peter-mount/golib/kernel"
	"github.com/peter-mount/golib/kernel/cron"
	"github.com/peter-mount/golib/rest"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwind3"
	"github.com/peter-mount/nre-feeds/darwindb"
)

type DarwinDBService struct {
	config      *bin.Config
	cron        *cron.CronService
	db          darwindb.DarwinDB
	em          *darwind3.DarwinEventManager
	restService *rest.Server
}

func (a *DarwinDBService) Name() string {
	return "DarwinDBService"
}

func (a *DarwinDBService) Init(k *kernel.Kernel) error {
	service, err := k.AddService(&bin.Config{})
	if err != nil {
		return err
	}
	a.config = (service).(*bin.Config)

	service, err = k.AddService(&cron.CronService{})
	if err != nil {
		return err
	}
	a.cron = (service).(*cron.CronService)

	service, err = k.AddService(&rest.Server{})
	if err != nil {
		return err
	}
	a.restService = (service).(*rest.Server)

	return nil
}

func (a *DarwinDBService) PostInit() error {
	err := a.db.Init(a.config)
	if err != nil {
		return err
	}

	return nil
}

func (a *DarwinDBService) Start() error {

	// Connect to Rabbit & name the connection so its easier to debug
	a.config.RabbitMQ.ConnectionName = "darwin db"
	err := a.config.RabbitMQ.Connect()
	if err != nil {
		return err
	}

	a.em = darwind3.NewDarwinEventManager(&a.config.RabbitMQ, a.config.D3.EventKeyPrefix)

	// Listen for deactivation messages
	err = a.em.RawListenToEvents(darwind3.Event_Deactivated, a.db.Deactivated)
	if err != nil {
		return err
	}

	// Schedule updates
	err = a.em.RawListenToEvents(darwind3.Event_ScheduleUpdated, a.db.ScheduleUpdated)
	if err != nil {
		return err
	}

	return nil
}

func (a *DarwinDBService) Stop() {
	a.db.Stop()
}
