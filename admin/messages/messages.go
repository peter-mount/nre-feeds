package messages

import (
  "github.com/peter-mount/go-kernel"
  "github.com/peter-mount/nre-feeds/admin"
  "github.com/peter-mount/nre-feeds/util/worker"
)

// Messages manages cleaning up station messages
type Messages struct {
  taskQueue *worker.TaskQueue
}

func (m *Messages) Name() string {
  return "AdminMessages"
}

func (m *Messages) Init(k *kernel.Kernel) error {
  service, err := k.AddService(&worker.TaskQueue{})
  if err != nil {
    return err
  }
  m.taskQueue = service.(*worker.TaskQueue)

  return k.DependsOn(&admin.Admin{})
}

func (m *Messages) Start() error {
  m.taskQueue.AddTask(&getStationMessages{})
  return nil
}
