package darwind3

import (
  "github.com/peter-mount/golib/rest"
  "log"
)

func (d *DarwinD3) RegisterRest( c *rest.ServerContext ) {
  c.Handle( "/message/broadcast", d.BroadcastStationMessagesHandler ).Methods( "POST" )
  c.Handle( "/message/{id}", d.StationMessageHandler ).Methods( "GET" )

  c.Handle( "/schedule/{rid}", d.ScheduleHandler ).Methods( "GET" )

  c.Handle( "/test", d.TestHandler ).Methods( "POST" )
}

// Test handle used to test xml locally via rest
func (d *DarwinD3) TestHandler( r *rest.Rest ) error {
  p := &Pport{}

  if err := r.Body( p ); err != nil {
    log.Println( err )
    return err
  }

  if err := p.Process( d ); err != nil {
    log.Println( err )
    return err
  }

  r.Status( 200 ).
    Value( "ok" )
  return nil
}
