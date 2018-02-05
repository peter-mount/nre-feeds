package darwind3

import (
  "github.com/peter-mount/golib/rest"
  "strconv"
)

// StationMessageHandler implements the /live/message/{id} rest endpoint
func (d *DarwinD3) StationMessageHandler( r *rest.Rest ) error {
  if id, err := strconv.Atoi( r.Var( "id" ) ); err != nil {
    r.Status( 404 )
  } else if msg := d.Messages.Get( id ); msg != nil {
    msg.Self = r.Self( r.Context() + "/message/" + r.Var( "id" ) )
    r.Status( 200 ).Value( msg )
  } else {
    r.Status( 404 )
  }

  return nil
}

// BroadcastStationMessagesHandler allows us to re-broadcast all messages
func (d *DarwinD3) BroadcastStationMessagesHandler( r *rest.Rest ) error {
  d.BroadcastStationMessages()
  r.Status( 200 ).Value( "OK" )

  return nil
}
