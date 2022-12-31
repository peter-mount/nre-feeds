package service

import (
	"github.com/gorilla/handlers"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwinkb"
)

type DarwinKBService struct {
	darwinkb *darwinkb.DarwinKB

	config *bin.Config

	restService *rest.Server
}

func (a *DarwinKBService) Name() string {
	return "DarwinKBService"
}

func (a *DarwinKBService) Init(k *kernel.Kernel) error {

	service, err := k.AddService(&bin.Config{})
	if err != nil {
		return err
	}
	a.config = (service).(*bin.Config)

	service, err = k.AddService(&darwinkb.DarwinKB{})
	if err != nil {
		return err
	}
	a.darwinkb = (service).(*darwinkb.DarwinKB)

	service, err = k.AddService(&rest.Server{})
	if err != nil {
		return err
	}
	a.restService = (service).(*rest.Server)

	// ReferenceUpdate
	return nil
}

func (a *DarwinKBService) Start() error {

	// nre-feeds#24 Add compression to output
	a.restService.Use(handlers.CompressHandler)

	a.restService.Handle("/companies", a.CompaniesHandler).Methods("GET")
	a.restService.Handle("/company/{id}", a.CompanyHandler).Methods("GET")

	a.restService.Handle("/incidents", a.IncidentsHandler).Methods("GET")
	a.restService.Handle("/incidents/{toc}", a.IncidentsTocHandler).Methods("GET")
	a.restService.Handle("/incident/{id}", a.IncidentHandler).Methods("GET")

	a.restService.Handle("/serviceIndicators", a.GetServiceIndicatorsHandler).Methods("GET")
	a.restService.Handle("/serviceIndicator/{id}", a.GetServiceIndicatorHandler).Methods("GET")

	a.restService.Handle("/station/{crs}", a.StationHandler).Methods("GET")

	a.restService.Handle("/ticket/types", a.TicketTypesHandler).Methods("GET")
	a.restService.Handle("/ticket/type/ids", a.TicketIdsHandler).Methods("GET")
	a.restService.Handle("/ticket/type/{id}", a.TicketTypeHandler).Methods("GET")

	// Expose the static directory so we offer the raw xml & full json files
	a.restService.Static("/static/", a.config.Database.KB+"static/")

	return nil
}
