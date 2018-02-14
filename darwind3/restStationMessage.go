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

// CrsMessageHandler Returns all messages for a CRS
func (d *DarwinD3) AllMessageHandler( r *rest.Rest ) error {
  var messages []*StationMessage

  d.Messages.Update( func() error {
    if len( d.Messages.messages ) > 0 {
      for _, s := range d.Messages.messages {
        messages = append( messages, s )
      }
    }
    return nil
  })

  r.Status( 200 ).Value( messages )

  return nil
}

// CrsMessageHandler Returns all messages for a CRS
func (d *DarwinD3) CrsMessageHandler( r *rest.Rest ) error {
  crs := r.Var( "crs" )

  var messages []*StationMessage

  d.Messages.Update( func() error {
    if len( d.Messages.messages ) > 0 {
      for _, s := range d.Messages.messages {
        for _, c := range s.Station {
          if c == crs {
            messages = append( messages, s )
            break
          }
        }
      }
    }
    return nil
  })

  r.Status( 200 ).Value( messages )

  return nil
}
