package worker

import (
  "context"
  "log"
)

const (
  queueKey = "__TaskQueue__"
)

type Task interface {
  Run(ctx context.Context) error
}

type TaskName interface {
  Name() string
}

type TaskQueue struct {
  ctx       context.Context
  tasks     []Task // Tasks to perform
  taskCount int    // Number of tasks, this is the number added not len(tasks)
}

func (a *TaskQueue) Name() string {
  return queueKey
}

func (a *TaskQueue) Start() error {
  ctx := context.Background()
  ctx = context.WithValue(ctx, queueKey, a)

  a.ctx = ctx
  return nil
}

func (a *TaskQueue) SetContext(key, value interface{}) *TaskQueue {
  a.ctx = context.WithValue(a.ctx, key, value)
  return a
}

func (a *TaskQueue) AddTask(t Task) {
  if t != nil {
    a.tasks = append(a.tasks, t)
    a.taskCount++
  }
}

func (a *TaskQueue) NextTask() Task {
  if len(a.tasks) == 0 {
    return nil
  }

  t := a.tasks[0]
  a.tasks = a.tasks[1:]
  return t
}

func AddTask(ctx context.Context, tasks ...Task) {
  queue := ctx.Value(queueKey).(*TaskQueue)
  if queue != nil {
    for _, task := range tasks {
      queue.AddTask(task)
    }
  }
}

func (a *TaskQueue) Run() error {
  tasksRun := 0
  defer func() {
    log.Printf("Performed %d/%d tasks", tasksRun, a.taskCount)
  }()

  for t := a.NextTask(); t != nil; t = a.NextTask() {

    // Optional TaskName interface
    if n, ok := t.(TaskName); ok {
      log.Printf("%03d %s", tasksRun+1, n.Name())
    }

    if err := t.Run(a.ctx); err != nil {
      return err
    }

    // only increment the run counter after a successful run
    tasksRun++
  }

  return nil
}
