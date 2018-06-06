package darwinref

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/rest"
  "log"
)

// viaHandler returns the unique instance of a via entry
func (dr *DarwinReference) ViaHandler( r *rest.Rest ) error {
  return dr.View( func( tx *bolt.Tx ) error {
    log.Printf( "via '%s' '%s' '%s' '%s'", r.Var( "at" ), r.Var( "dest" ), r.Var( "loc1" ), r.Var( "loc2" ) )

    if via, exists := dr.GetVia( tx, r.Var( "at" ), r.Var( "dest" ), r.Var( "loc1" ), r.Var( "loc2" ) ); exists {
      via.SetSelf( r )
      r.Status( 200 ).Value( via )
    } else {
      r.Status( 404 )
    }

    return nil
  })
}
