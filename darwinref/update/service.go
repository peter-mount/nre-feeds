package update

import (
	"github.com/peter-mount/golib/kernel"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwind3"
	"github.com/peter-mount/nre-feeds/darwinref/service"
)

type ReferenceUpdateService struct {
	config *bin.Config
	ref    *service.DarwinRefService
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

	return nil
}

func (a *ReferenceUpdateService) Start() error {
	// Only listen if S3 is enabled
	if a.config.S3.Enabled {
		a.config.RabbitMQ.ConnectionName = "darwin ref"
		a.config.RabbitMQ.Connect()

		em := darwind3.NewDarwinEventManager(&a.config.RabbitMQ, a.config.D3.EventKeyPrefix)
		em.ListenToEvents(darwind3.Event_TimeTableUpdate, a.referenceUpdateListener)
	}

	return nil
}
