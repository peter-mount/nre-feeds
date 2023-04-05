package darwind3

import (
	"encoding/json"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/rabbitmq"
	"log"
	"os"
)

// The possible types of DarwinEvent
const (
	// A schedule was deactivated
	Event_Deactivated = "deactivated"
	// A schedule was updated
	Event_ScheduleUpdated = "scheduleUpdated"
	// A new StationMessage
	Event_StationMessage = "stationMessage"
	// A station's departure boards have been updated (LDB only)
	Event_BoardUpdate = "boardUpdate"
	// TimeTable update (either timetable or reference)
	Event_TimeTableUpdate = "timeTableUpdate"
	// TrackingID update
	Event_TrackingID = "trackingID"
	// Alarms
	Event_Alarm = "alarm"
	// Request station message broadcast
	Event_Request_StationMessage = "requestStationMessage"
)

// An event notifying of something happening within DarwinD3
type DarwinEvent struct {
	// The type of the event
	Type string
	// The RID of the train that caused this event
	RID string
	// The affected schedule or nil if none
	Schedule *Schedule
	// The CRS code of the station in this event (LDB only)
	Crs string
	// The StationMessage that's been updated
	NewStationMessage *StationMessage
	// The existing message before the update (or nil)
	ExistingStationMessage *StationMessage
	// TimeTable update
	TimeTableId *TimeTableId
	// TrackingID update
	TrackingID *TrackingID
	// Alarm
	Alarm   *Alarm
	AlarmId string
}

// The core of the eventing system
type DarwinEventManager struct {
	mq             *rabbitmq.RabbitMQ
	prefix         string
	eventKeyPrefix string
}

// NewDarwinEventManager creates a new DarwinEventManager
func NewDarwinEventManager(mq *rabbitmq.RabbitMQ, eventKeyPrefix string) *DarwinEventManager {
	d := &DarwinEventManager{}
	d.mq = mq

	// Queue prefix, try to use the local hostname (e.g. of the container)
	if hostname, err := os.Hostname(); err != nil {
		d.prefix = "error"
	} else {
		d.prefix = hostname
	}

	// The eventKeyPrefix added to the routingKey & queueName to keep them unique
	if eventKeyPrefix != "" {
		d.eventKeyPrefix = eventKeyPrefix + "."
	}

	return d
}

func (d *DarwinEventManager) RawListenToEvents(eventType string, f func([]byte)) error {
	queueName := fmt.Sprintf("%s.%sd3.event.%s", d.prefix, d.eventKeyPrefix, eventType)
	routingKey := fmt.Sprintf("%sd3.event.%s", d.eventKeyPrefix, eventType)

	if channel, err := d.mq.NewChannel(); err != nil {
		log.Println(err)
		return err
	} else {

		// Force prefetchCount to 1 so we don't get everything in one go
		_ = channel.Qos(1, 0, false)

		// non-durable auto-delete queue
		_, _ = d.mq.QueueDeclare(channel, queueName, false, true, false, false, nil)

		_ = d.mq.QueueBind(channel, queueName, routingKey, "amq.topic", false, nil)

		ch, _ := d.mq.Consume(channel, queueName, "D3 Event Consumer", true, true, false, false, nil)

		go func() {
			for {
				msg := <-ch
				f(msg.Body)
			}
		}()

		return nil
	}
}

// ListenToEvents will run a function which will reveive DarwinEvent's for the
// specified type until it exists.
func (d *DarwinEventManager) ListenToEvents(eventType string, f func(*DarwinEvent)) error {
	return d.RawListenToEvents(eventType, func(body []byte) {
		evt := &DarwinEvent{}

		_ = json.Unmarshal(body, evt)

		if evt.Type == eventType {
			if evt.Schedule != nil {
				evt.Schedule.Sort()
			}
			f(evt)
		}
	})
}

// PostEvent posts a DarwinEvent to all listeners listening for that specific type
func (d *DarwinEventManager) PostEvent(e *DarwinEvent) {
	if b, err := json.Marshal(e); err == nil {
		_ = d.mq.Publish(fmt.Sprintf("%sd3.event.%s", d.eventKeyPrefix, e.Type), b)
	}
}
