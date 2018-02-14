package darwind3

import (
  "github.com/peter-mount/golib/rest"
)

func (d *DarwinD3) RegisterRest( c *rest.ServerContext ) {
  c.Handle( "/message/broadcast", d.BroadcastStationMessagesHandler ).Methods( "POST" )

  c.Handle( "/messages", d.AllMessageHandler ).Methods( "GET" )
  c.Handle( "/messages/{crs}", d.CrsMessageHandler ).Methods( "GET" )

  c.Handle( "/message/{id}", d.StationMessageHandler ).Methods( "GET" )

  c.Handle( "/schedule/{rid}", d.ScheduleHandler ).Methods( "GET" )
}
