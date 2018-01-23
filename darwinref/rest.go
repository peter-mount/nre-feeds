package darwinref

import (
  "github.com/peter-mount/golib/rest"
)

// RegisterRest registers the rest endpoints into a ServerContext
func (r DarwinReference) RegisterRest( c *rest.ServerContext ) {
  c.Handle( "/import", r.ImportHandler ).Methods( "POST" )
}
