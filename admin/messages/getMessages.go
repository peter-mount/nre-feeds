package messages

import (
  "context"
  "github.com/peter-mount/nre-feeds/darwind3/client"
  "github.com/peter-mount/nre-feeds/util/worker"
  "time"
)

type getStationMessages struct {
}

func (m *getStationMessages) Name() string {
  return "Retrieve all active station messages"
}

func (m *getStationMessages) Run(ctx context.Context) error {
  d3Client := ctx.Value("d3").(*client.DarwinD3Client)

  messages, err := d3Client.GetStationMessages()
  if err != nil {
    return err
  }

  limit := time.Now().Add(-24 * time.Hour)

  for _, m := range messages {
    if m.Date.Before(limit) {
      worker.AddTask(ctx, &pruneStationMessage{ID: m.ID})
    }
  }
  return nil
}
