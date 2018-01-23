package darwinref

import (
  "github.com/peter-mount/golib/rest"
  "log"
)

func (dr *DarwinReference) ImportHandler( r *rest.Rest ) error {
  log.Println( "DarwinReference import: started" )

  if err := r.Body( dr ); err != nil {
    return err
  }

  log.Println( "DarwinReference import: completed" )
  r.Status( 200 ).
    Value( "ok" )
    return nil
}
