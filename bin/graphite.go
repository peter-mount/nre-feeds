package bin

import (
  "fmt"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/statistics"
  "github.com/streadway/amqp"
  "strings"
)

type Graphite struct {
  statistics    statistics.Statistics
  channel      *amqp.Channel
  config       *Config
}

func (g *Graphite) Name() string {
  return "Graphite"
}

func (a *Graphite) Init( k *kernel.Kernel ) error {
  service, err := k.AddService( &Config{} )
  if err != nil {
    return err
  }
  a.config = (service).(*Config)

  return nil
}

func (g *Graphite) Start() error {
  // Custom statistics engine, capture every 10s so we submit to Graphite at
  // intervals it's expecting
  g.statistics.Log = false
  g.statistics.Schedule = "0/10 * * * * *"
  g.statistics.Configure()

  if g.config.Graphite.Enabled  {

    // Default exchange is "graphite"
    if g.config.Graphite.Exchange == "" {
      g.config.Graphite.Exchange = "graphite"
    }

    err := g.config.RabbitMQ.Connect()
    if err != nil {
      return err
    }

    g.channel, err = g.config.RabbitMQ.NewChannel()
    if err != nil {
      return err
    }

    // We are a statistics Recorder
    g.statistics.Recorder = g
  }

  return nil
}

// PublishStatistic Handles publishing statistics to Graphite over RabbitMQ
func (g *Graphite) PublishStatistic( name string, s *statistics.Statistic ) {
  if strings.HasSuffix( name, "td.all" ) {
    // Value will be the latency
    g.publish( name + ".latency", s.Value, s.Timestamp )
    // Count the number of messages
    g.publish( name + ".count", s.Count, s.Timestamp )

    // Min/Max latency values - don't send if max<min - i.e. no data!
    if s.Max >= s.Min {
      g.publish( name + ".min", s.Min, s.Timestamp )
      g.publish( name + ".max", s.Max, s.Timestamp )
    }
  } else {
    g.publish( name, s.Value, s.Timestamp )
  }
}

func (g *Graphite) publish( name string, val int64, ts int64 ) {
  statName := name
  if g.config.Graphite.Prefix != "" {
    statName = g.config.Graphite.Prefix + "." + name
  }
  msg := fmt.Sprintf( "%s %d %d", statName, val, ts)

  g.channel.Publish(
    g.config.Graphite.Exchange,
    statName,
    false,
    false,
    amqp.Publishing{
      Body: []byte(msg),
  })
}
