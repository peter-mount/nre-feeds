package service

import (
	"github.com/gorilla/handlers"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/go-kernel/v2/cron"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwind3"
	"github.com/peter-mount/nre-feeds/ldb"
	"runtime/debug"
)

type LDBService struct {
	ldb         ldb.LDB
	config      *bin.Config
	cron        *cron.CronService
	restService *rest.Server
}

func (a *LDBService) Name() string {
	return "LDBService"
}

func (a *LDBService) Init(k *kernel.Kernel) error {
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

	// ReferenceUpdate
	return nil
}

func (a *LDBService) PostInit() error {
	a.ldb.Darwin = a.config.Services.DarwinD3
	a.ldb.Reference = a.config.Services.Reference

	// nre-feeds#24 Add compression to output
	a.restService.Use(handlers.CompressHandler)

	// Rest services
	a.restService.Handle("/boards/{crs}", a.stationHandler).Methods("GET")
	a.restService.Handle("/service/{rid}", a.serviceHandler).Methods("GET")

	return nil
}

func (a *LDBService) Start() error {

	// Connect to Rabbit & name the connection so its easier to debug
	a.config.RabbitMQ.ConnectionName = "darwin ldb"
	err := a.config.RabbitMQ.Connect()
	if err != nil {
		return err
	}

	a.ldb.EventManager = darwind3.NewDarwinEventManager(&a.config.RabbitMQ, a.config.D3.EventKeyPrefix)

	a.config.DbPath(&a.config.Database.LDB, "ldb.db")

	err = a.ldb.Init(a.config.Database.LDB)
	if err != nil {
		return err
	}

	// Purge old schedules every hour
	_, err = a.cron.AddFunc("0 5 * * * *", a.ldb.PurgeSchedules)
	if err != nil {
		return err
	}

	// Check for any orphans once every 6 hours
	_, err = a.cron.AddFunc("0 10 0/6 * * *", a.ldb.PurgeOrphans)
	if err != nil {
		return err
	}

	// Check for any orphaned services every 6 hours
	_, err = a.cron.AddFunc("0 15 0/6 * * *", a.ldb.PurgeServices)
	if err != nil {
		return err
	}

	// Log DB status every hour
	_, err = a.cron.AddFunc("0 10 * * * *", a.ldb.DBStatus)
	if err != nil {
		return err
	}

	// Memory
	_, err = a.cron.AddFunc("9/10 * * * * *", func() {
		darwind3.SubmitMemStats("darwin.ldb")
	})
	if err != nil {
		return err
	}

	_, err = a.cron.AddFunc("0 0/5 * * * *", debug.FreeOSMemory)
	if err != nil {
		return err
	}

	return nil
}
