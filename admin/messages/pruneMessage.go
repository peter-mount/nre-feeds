package messages

import (
  "context"
  "fmt"
  "github.com/peter-mount/go-kernel/rabbitmq"
  "time"
)

type pruneStationMessage struct {
  ID int64
}

func (m *pruneStationMessage) Name() string {
  return fmt.Sprintf("Remove Station Message %d", m.ID)
}

const (
  TIMESTAMP = "2006-01-02T15:04:05Z"
)

func (m *pruneStationMessage) Run(ctx context.Context) error {
  //d3Client := ctx.Value("d3").(*client.DarwinD3Client)
  mq := ctx.Value("mq").(*rabbitmq.RabbitMQ)

  mq.Publish(
    "nre.darwin.pushport-v16",
    []byte(fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>"+
        "<Pport xmlns=\"http://www.thalesgroup.com/rtti/PushPort/v16\" xmlns:ns7=\"http://www.thalesgroup.com/rtti/PushPort/StationMessages/v1\" ts=\"%s\" version=\"16.0\">"+
        "<sR><OW id=\"%d\" cat=\"Train\" sev=\"1\"><ns7:Msg>Deleted</ns7:Msg></OW></sR>"+
        "</Pport>",
      time.Now().UTC().Format(TIMESTAMP),
      m.ID)))

  return nil
}
