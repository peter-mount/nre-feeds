package update

import (
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwind3"
	"github.com/peter-mount/nre-feeds/darwintimetable/service"
)

type TimetableUpdateService struct {
	config    *bin.Config
	timetable *service.DarwinTimetableService
}

func (a *TimetableUpdateService) Name() string {
	return "TimetableUpdateService"
}

func (a *TimetableUpdateService) Init(k *kernel.Kernel) error {
	svce, err := k.AddService(&bin.Config{})
	if err != nil {
		return err
	}
	a.config = (svce).(*bin.Config)

	svce, err = k.AddService(&service.DarwinTimetableService{})
	if err != nil {
		return err
	}
	a.timetable = (svce).(*service.DarwinTimetableService)

	return nil
}

func (a *TimetableUpdateService) Start() error {
	// Only listen if S3 is enabled
	if a.config.S3.Enabled {
		a.config.RabbitMQ.ConnectionName = "darwin tt"
		err := a.config.RabbitMQ.Connect()
		if err != nil {
			return err
		}

		em := darwind3.NewDarwinEventManager(&a.config.RabbitMQ, a.config.D3.EventKeyPrefix)
		err = em.ListenToEvents(darwind3.Event_TimeTableUpdate, a.timetableUpdateListener)
		if err != nil {
			return err
		}

		go a.findUpdates()
	}

	return nil
}
