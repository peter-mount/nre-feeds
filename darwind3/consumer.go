package darwind3

import (
  "bytes"
  "encoding/xml"
  "github.com/peter-mount/golib/statistics"
  "github.com/peter-mount/golib/rabbitmq"
  "github.com/streadway/amqp"
)

// BindConsumer binds a consumer to a RabbitMQ queue to receive D3 messages
func (d *DarwinD3) BindConsumer( r *rabbitmq.RabbitMQ, queueName, routingKey string ) error {
  if channel, err := r.NewChannel(); err != nil {
    return err
  } else {
    r.QueueDeclare( channel, queueName, true, false, false, false, nil )
    r.QueueBind( channel, queueName, routingKey, "amq.topic", false, nil )
    ch, _ := r.Consume( channel, queueName, "ldb consumer", false, true, false, false, nil )

    go func() {
      for {
        msg := <- ch
        d.consume( msg )
      }
    }()

    return nil
  }
}

func (d *DarwinD3) consume( msg amqp.Delivery ) {
  defer msg.Ack( false )

  reader := bytes.NewReader( msg.Body )
  p := &Pport{}
  if err := xml.NewDecoder( reader ).Decode( p ); err == nil {
    p.Process( d )
    statistics.Incr( "d3.in" )
  }
}
