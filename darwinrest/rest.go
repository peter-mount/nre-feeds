// darwinrest provides some additional rest services which use all of the other
// packages in forming their results
package darwinrest

import (
  "darwinref"
  "darwintimetable"
  "github.com/peter-mount/golib/rest"
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
