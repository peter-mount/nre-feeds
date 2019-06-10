package service

import (
	"github.com/peter-mount/filecache"
	fcsve "github.com/peter-mount/filecache/service"
	"github.com/peter-mount/golib/kernel"
	"github.com/peter-mount/golib/kernel/cron"
	"github.com/peter-mount/golib/rest"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwind3"
	"runtime/debug"
	"time"
)

type DarwinD3Service struct {
	darwind3    darwind3.DarwinD3
	config      *bin.Config
	cron        *cron.CronService
	restService *rest.Server
	fileCache   *fcsve.FileCacheService
}

const (
	metaExpiry            = 24 * time.Hour
	stationMessagesExpiry = 12 * time.Hour
	scheduleDiskExpiry    = 4 * 24 * time.Hour
)

func (a *DarwinD3Service) Name() string {
	return "DarwinD3Service"
}

func (a *DarwinD3Service) Init(k *kernel.Kernel) error {
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

	service, err = k.AddService(&fcsve.FileCacheService{})
	if err != nil {
		return err
	}
	a.fileCache = (service).(*fcsve.FileCacheService)

	return nil
}

func (a *DarwinD3Service) PostInit() error {
	a.darwind3.Config = a.config

	if a.config.D3.ResolveSched {
		// Allow D3 to resolve schedules from the timetable
		a.darwind3.Timetable = a.config.Services.Timetable
	}

	// Rest services
	a.restService.Handle("/alarm/{rid}", a.AlarmHandler).Methods("GET")
	a.restService.Handle("/alarms", a.AlarmsHandler).Methods("GET")

	a.restService.Handle("/message/broadcast", a.BroadcastStationMessagesHandler).Methods("POST")
	a.restService.Handle("/message/{id}", a.StationMessageHandler).Methods("GET")
	a.restService.Handle("/messages", a.AllMessageHandler).Methods("GET")
	a.restService.Handle("/messages/{crs}", a.CrsMessageHandler).Methods("GET")

	a.restService.Handle("/schedule/{rid}", a.ScheduleHandler).Methods("GET")

	a.restService.Handle("/status", a.StatusHandler).Methods("GET")

	var err error
	cache := a.fileCache.Cache()

	a.darwind3.Alarms, err = cache.AddCache(filecache.CacheTableConfig{
		Name:           "alarms",
		ExpiryTime:     time.Hour,
		FromBytes:      darwind3.AlarmFromBytes,
		ToBytes:        filecache.ToJsonBytes,
		StartupOptions: filecache.ExpireCacheOnStart,
		DiskExpiryTime: metaExpiry,
	})
	if err != nil {
		return err
	}

	a.darwind3.Associations, err = cache.AddCache(filecache.CacheTableConfig{
		Name:           "associations",
		ExpiryTime:     10 * time.Minute,
		FromBytes:      darwind3.AssociationsFromBytes,
		ToBytes:        filecache.ToJsonBytes,
		StartupOptions: filecache.ExpireCacheOnStart,
		DiskExpiryTime: scheduleDiskExpiry,
	})
	if err != nil {
		return err
	}

	a.darwind3.Meta, err = cache.AddCache(filecache.CacheTableConfig{
		Name:           "meta",
		ExpiryTime:     24 * time.Hour,
		FromBytes:      filecache.TimeFromBytes,
		ToBytes:        filecache.ToJsonBytes,
		StartupOptions: filecache.ExpireCacheOnStart,
		DiskExpiryTime: metaExpiry,
	})
	if err != nil {
		return err
	}

	a.darwind3.Schedules, err = cache.AddCache(filecache.CacheTableConfig{
		Name:           "schedules",
		ExpiryTime:     120 * time.Second,
		FromBytes:      darwind3.FromBytesSchedule,
		ToBytes:        filecache.ToJsonBytes,
		StartupOptions: filecache.ExpireCacheOnStart,
		DiskExpiryTime: scheduleDiskExpiry,
	})
	if err != nil {
		return err
	}

	a.darwind3.StationMessages, err = cache.AddCache(filecache.CacheTableConfig{
		Name:           "stationMessages",
		ExpiryTime:     24 * time.Hour,
		FromBytes:      darwind3.StationMessageFromBytes,
		ToBytes:        filecache.ToJsonBytes,
		StartupOptions: filecache.ExpireCacheOnStart,
		DiskExpiryTime: stationMessagesExpiry,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *DarwinD3Service) Start() error {

	// Connect to Rabbit & name the connection so its easier to debug
	a.config.RabbitMQ.ConnectionName = "darwin d3"
	err := a.config.RabbitMQ.Connect()
	if err != nil {
		return err
	}

	em := darwind3.NewDarwinEventManager(&a.config.RabbitMQ, a.config.D3.EventKeyPrefix)

	a.config.DbPath(&a.config.Database.PushPort, "dwd3.db")
	a.darwind3.Init(em)

	// Memory
	_, err = a.cron.AddFunc("9/10 * * * * *", func() {
		darwind3.SubmitMemStats("darwin.d3")
	})
	if err != nil {
		return err
	}

	_, err = a.cron.AddFunc("0 0/5 * * * *", debug.FreeOSMemory)
	if err != nil {
		return err
	}
	/*
	   _, err = a.cron.AddFunc("0 * * * * *", darwind3.GC)
	   if err != nil {
	     return err
	   }
	*/

	// Listen for broadcast events
	err = a.darwind3.EventManager.ListenToEvents(darwind3.Event_Request_StationMessage, a.darwind3.BroadcastStationMessages)
	if err != nil {
		return err
	}

	// The V16 PushPort queue
	err = a.darwind3.BindConsumer(&a.config.RabbitMQ, a.config.D3.PushPort.QueueName, a.config.D3.PushPort.RoutingKey)
	if err != nil {
		return err
	}

	return nil
}
