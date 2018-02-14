package darwind3

import (
  "encoding/json"
  "fmt"
  "github.com/peter-mount/golib/rabbitmq"
  "os"
)

// The possible types of DarwinEvent
const (
  // A schedule was deactivated
  Event_Deactivated = iota
  // A schedule was updated
  Event_ScheduleUpdated
  // A new StationMessage
  Event_StationMessage
  // A station's departure boards have been updated (LDB only)
  Event_BoardUpdate
)

// An event notifying of something happening within DarwinD3
type DarwinEvent struct {
  // The type of the event
  Type        int
  // The RID of the train that caused this event
  RID         string
  // The affected Schedule or nil if none
  Schedule   *Schedule
  // The CRS code of the station in this event (LDB only)
  Crs         string
  // The StationMessage that's been updated
  NewStationMessage        *StationMessage
  // The existing message before the update (or nil)
  ExistingStationMessage   *StationMessage
}

// The core of the eventing system
type DarwinEventManager struct {
  mq         *rabbitmq.RabbitMQ
  prefix      string
  sequence    int
}

// NewDarwinEventManager creates a new DarwinEventManager
func NewDarwinEventManager( mq *rabbitmq.RabbitMQ ) *DarwinEventManager {
  d := &DarwinEventManager{}
  d.mq = mq
  if hostname, err := os.Hostname(); err != nil {
    d.prefix = "error"
  } else {
    d.prefix = hostname
  }
  return d
}

// ListenToEvents will run a function which will reveive DarwinEvent's for the
// specified type until it exists.
func (d *DarwinEventManager) ListenToEvents( eventType int, f func( *DarwinEvent ) ) {
  d.mq.Connect()
  seq := d.sequence
  d.sequence++

  queueName := fmt.Sprintf( "%s.d3.event.%d.%d", d.prefix, eventType, seq)
  routingKey := fmt.Sprintf( "d3.event.%d", eventType )

  // non-durable auto-delete queue
  d.mq.QueueDeclare( queueName, false, true, false, false, nil )

  d.mq.QueueBind( queueName, routingKey, "amq.topic", false, nil )

  ch, _ := d.mq.Consume( queueName, "D3 Event Consumer", true, true, false, false, nil )

  go func() {
    for {
      msg := <- ch

      evt := &DarwinEvent{}
      json.Unmarshal( msg.Body, evt )

      if evt.Type == eventType {
        f( evt )
      }
    }
  }()

}

// PostEvent posts a DarwinEvent to all listeners listening for that specific type
func (d *DarwinEventManager) PostEvent( e *DarwinEvent ) {
  if b, err := json.Marshal( e ); err == nil {
    d.mq.Publish( fmt.Sprintf( "d3.event.%d", e.Type ), b )
  }
}
