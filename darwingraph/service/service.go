package service

import (
	"github.com/gorilla/handlers"
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/go-kernel/cron"
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwingraph"
)

type DarwinGraphService struct {
	config      *bin.Config
	cron        *cron.CronService
	restService *rest.Server
	darwinGraph *darwingraph.DarwinGraph
}

func (a *DarwinGraphService) Name() string {
	return "DarwinGraphService"
}

func (a *DarwinGraphService) Init(k *kernel.Kernel) error {
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

	service, err = k.AddService(&darwingraph.DarwinGraph{})
	if err != nil {
		return err
	}
	a.darwinGraph = (service).(*darwingraph.DarwinGraph)
	return nil
}

func (a *DarwinGraphService) PostInit() error {
	/*
		a.config.DbPath(&a.config.Database.Reference, "dwref.db")
		if err := a.reference.OpenDB(a.config.Database.Reference); err != nil {
			return err
		}
	*/
	// nre-feeds#24 Add compression to output
	a.restService.Use(handlers.CompressHandler)

	// Rest services
	/*
		a.restService.Handle("/reason/cancelled", a.AllReasonCancelHandler).Methods("GET")
		a.restService.Handle("/reason/cancelled/{id}", a.ReasonCancelHandler).Methods("GET")

		a.restService.Handle("/reason/late", a.AllReasonLateHandler).Methods("GET")
		a.restService.Handle("/reason/late/{id}", a.ReasonLateHandler).Methods("GET")

		a.restService.Handle("/via/{at}/{dest}/{loc1}", a.ViaHandler).Methods("GET")
		a.restService.Handle("/via/{at}/{dest}/{loc1}/{loc2}", a.ViaHandler).Methods("GET")
		a.restService.Handle("/via", a.ViaResolveHandler).Methods("POST")

		// Reference retrieval methods
		a.restService.Handle("/crs/{id}", a.CrsHandler).Methods("GET")
		a.restService.Handle("/crs", a.AllCrsHandler).Methods("GET")

		a.restService.Handle("/tiploc", a.TiplocsHandler).Methods("POST")
		a.restService.Handle("/tiploc/{id}", a.TiplocHandler).Methods("GET")

		a.restService.Handle("/toc", a.AllTocsHandler).Methods("GET")
		a.restService.Handle("/toc/{id}", a.TocHandler).Methods("GET")

		a.restService.Handle("/search/{term}", a.SearchHandler).Methods("GET")
	*/
	return nil
}

func (a *DarwinGraphService) Stop() {
	//a.reference.Close()
}
