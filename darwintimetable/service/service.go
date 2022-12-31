package service

import (
	"github.com/gorilla/handlers"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/go-kernel/v2/cron"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwintimetable"
)

type DarwinTimetableService struct {
	timetable darwintimetable.DarwinTimetable

	config      *bin.Config
	cron        *cron.CronService
	restService *rest.Server
}

func (a *DarwinTimetableService) GetTimetable() *darwintimetable.DarwinTimetable {
	return &a.timetable
}

func (a *DarwinTimetableService) Name() string {
	return "DarwinRefService"
}

func (a *DarwinTimetableService) Init(k *kernel.Kernel) error {
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

func (a *DarwinTimetableService) PostInit() error {
	a.config.DbPath(&a.config.Database.Timetable, "dwtt.db")
	if err := a.timetable.OpenDB(a.config.Database.Timetable); err != nil {
		return err
	}

	// Prune schedules at 2am
	_, _ = a.cron.AddFunc("0 0 2 * * *", func() {
		a.timetable.PruneSchedules()
	})

	/*
			   if a.config.Ftp.Enabled {
			     // Scheduled updates
		       if a.config.Ftp.schedule != "" {
		         a.cron.AddFunc( a.config.Ftp.schedule, func () {
			         if err := a.updater.TimetableUpdate( &a.timetable ); err != nil {
			           log.Println( "Failed import:", err )
			         }
			       })
		         log.Println( "Auto SnapshotUpdate using:", a.config.Ftp.schedule )
			     }

			     // Initial import required?
			     if a.updater.ImportRequiredTimetable( &a.timetable ) {
			       a.updater.TimetableUpdate( &a.timetable )
			     }
			   }
	*/

	// nre-feeds#24 Add compression to output
	a.restService.Use(handlers.CompressHandler)

	// Rest services
	a.restService.Handle("/journey/{rid}", a.JourneyHandler).Methods("GET")

	return nil
}

func (a *DarwinTimetableService) Stop() {
	a.timetable.Close()
}
