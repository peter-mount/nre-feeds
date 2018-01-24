package darwind3

import (
  "github.com/peter-mount/golib/rest"
  "log"
)

func (d *DarwinD3) SetupRest( c *rest.ServerContext ) {
  c.Handle( "/test", d.TestHandler ).Methods( "POST" )
}

// Test handle used to test xml locally via rest
func (d *DarwinD3) TestHandler( r *rest.Rest ) error {
  pp := &Pport{}

  if err := r.Body( pp ); err != nil {
    log.Println( err )
    return err
  }

  //log.Printf( "Pport %v\n", pp )

  if err := pp.Process( d ); err != nil {
    log.Println( err )
    return err
  }

  r.Status( 200 ).
    Value( "ok" )
  return nil
}
