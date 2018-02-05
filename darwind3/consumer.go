package darwind3

import (
  "bytes"
  "encoding/xml"
  "github.com/peter-mount/golib/statistics"
  "github.com/peter-mount/golib/rabbitmq"
  "github.com/streadway/amqp"
)

// BindConsumer binds a consumer to a RabbitMQ queue to receive D3 messages
func (d *DarwinD3) BindConsumer( r *rabbitmq.RabbitMQ, queueName, routingKey string ) {
  r.Connect()
  r.QueueDeclare( queueName, true, false, false, false, nil )
  r.QueueBind( queueName, routingKey, "amq.topic", false, nil )
  ch, _ := r.Consume( queueName, "ldb consumer", false, true, false, false, nil )
  go func() {
    for {
      msg := <- ch
      d.consume( msg )
    }
  }()
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
