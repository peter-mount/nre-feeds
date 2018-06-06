// darwind3 Microservice
package main

import (
  "github.com/peter-mount/nre-feeds/bin"
  "github.com/peter-mount/nre-feeds/darwind3"
)

func main() {
  bin.RunApplication( app )
}

func app( config *bin.Config ) ( func(), error ) {

  // Connect to Rabbit & name the connection so its easier to debug
  config.RabbitMQ.ConnectionName = "darwin d3"
  config.RabbitMQ.Connect()

  d3 := &darwind3.DarwinD3{}

  config.DbPath( &config.Database.PushPort, "dwd3.db" )

  em := darwind3.NewDarwinEventManager( &config.RabbitMQ )

  if err := d3.OpenDB( config.Database.PushPort, em ); err != nil {
    return nil, err
  }

  if config.D3.ResolveSched {
    // Allow D3 to resolve schedules from the timetable
    d3.Timetable = config.Services.Timetable
  }

  d3.RegisterRest( config.Server.Ctx )

  // Expire old messages every 15 minutes & run an expire on startup
  config.Cron.AddFunc( "0 0/15 * * * *", d3.ExpireStationMessages )
  go d3.ExpireStationMessages()

  if config.RabbitMQ.Url != "" {
    d3.BindConsumer( &config.RabbitMQ, config.D3.QueueName, config.D3.RoutingKey )
  }

  return nil, nil
}
