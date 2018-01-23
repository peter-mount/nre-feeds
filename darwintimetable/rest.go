package darwintimetable

import (
  "github.com/peter-mount/golib/rest"
)

// RegisterRest registers the rest endpoints into a ServerContext
func (r DarwinTimetable) RegisterRest( c *rest.ServerContext ) {

  // Data import
  c.Handle( "/import", r.ImportHandler ).Methods( "POST" )
}
