package admin

import (
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwind3/client"
	"github.com/peter-mount/nre-feeds/util/worker"
)

type Admin struct {
	taskQueue *worker.TaskQueue
	config    *bin.Config
}

func (a *Admin) Name() string {
	return "Admin"
}

func (a *Admin) Init(k *kernel.Kernel) error {
	service, err := k.AddService(&bin.Config{})
	if err != nil {
		return err
	}
	a.config = service.(*bin.Config)

	service, err = k.AddService(&worker.TaskQueue{})
	if err != nil {
		return err
	}
	a.taskQueue = service.(*worker.TaskQueue)

	return err
}

func (a *Admin) Start() error {
	a.taskQueue.SetContext("d3", &client.DarwinD3Client{Url: a.config.Services.DarwinD3})

	a.config.RabbitMQ.ConnectionName = "darwin admin"

	if err := a.config.RabbitMQ.Connect(); err != nil {
		return err
	}

	a.taskQueue.SetContext("mq", &a.config.RabbitMQ)

	return nil
}
