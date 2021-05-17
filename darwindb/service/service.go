package service

import (
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/go-kernel/cron"
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwind3"
	"github.com/peter-mount/nre-feeds/darwindb"
)

type DarwinDBService struct {
	config      *bin.Config
	cron        *cron.CronService
	db          darwindb.DarwinDB
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

	a.restService.Handle("/service/{rid}", a.getService).Methods("GET")
	a.restService.Handle("/services/{crs}/{date}/{hour}", a.getStationServices).Methods("GET")

	return nil
}

func (a *DarwinDBService) Start() error {

	// Connect to Rabbit & name the connection so its easier to debug
	a.config.RabbitMQ.ConnectionName = "darwin db"
	err := a.config.RabbitMQ.Connect()
	if err != nil {
		return err
	}

	err = a.db.Subscribe(
		&a.config.RabbitMQ,
		a.config.D3.EventKeyPrefix,
		"schedules",
		a.db.ScheduleUpdated,
		darwind3.Event_ScheduleUpdated,
		// Don't include deactivated schedules as we already have them & this removes data
		//darwind3.Event_Deactivated,
	)
	if err != nil {
		return err
	}

	// Run the index job every 10 seconds
	_, err = a.cron.AddFunc("1/10 * * * * *", a.db.IndexSchedules)
	if err != nil {
		return err
	}

	return nil
}

func (a *DarwinDBService) Stop() {
	a.db.Stop()
}
