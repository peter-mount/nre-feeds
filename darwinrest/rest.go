// darwinrest provides some additional rest services which use all of the other
// packages in forming their results
package darwinrest

import (
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/nre-feeds/darwinref"
  "github.com/peter-mount/nre-feeds/darwintimetable"
)

type DarwinRest struct {
  Ref  *darwinref.DarwinReference
  TT   *darwintimetable.DarwinTimetable
}

// RegisterRest registers the rest endpoints into a ServerContext
func (r DarwinRest) RegisterRest( c *rest.ServerContext ) {

  if r.Ref != nil && r.TT != nil {
    c.Handle( "/journey/{rid}", r.JourneyHandler ).Methods( "GET" )
  }
}
