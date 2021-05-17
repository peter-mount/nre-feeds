package darwind3

import (
	"bytes"
	"encoding/xml"
	"github.com/peter-mount/go-kernel/rabbitmq"
	"github.com/peter-mount/golib/statistics"
	"github.com/streadway/amqp"
)

// BindConsumer binds a consumer to a RabbitMQ queue to receive D3 messages
func (d *DarwinD3) BindConsumer(r *rabbitmq.RabbitMQ, queueName, routingKey string) error {
	channel, err := r.NewChannel()
	if err != nil {
		return err
	}

	// Force prefetchCount to 1 so we don't get everything in one go
	_ = channel.Qos(1, 0, false)

	_, _ = r.QueueDeclare(channel, queueName, true, false, false, false, nil)
	_ = r.QueueBind(channel, queueName, routingKey, "amq.topic", false, nil)
	ch, _ := r.Consume(channel, queueName, "ldb consumer", false, true, false, false, nil)

	go func() {
		for {
			msg := <-ch
			d.consume(msg)
		}
	}()

	return nil
}

func (d *DarwinD3) consume(msg amqp.Delivery) {
	defer msg.Ack(false)

	p := &Pport{}

	p.FeedHeaders.populate(msg)

	reader := bytes.NewReader(msg.Body)
	if err := xml.NewDecoder(reader).Decode(p); err == nil {
		d.FeedStatus.process(p)
		_ = p.Process(d)
		statistics.Incr("d3.in")
	}
}
