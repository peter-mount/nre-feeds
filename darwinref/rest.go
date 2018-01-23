package darwinref

import (
  "github.com/peter-mount/golib/rest"
)

// RegisterRest registers the rest endpoints into a ServerContext
func (r DarwinReference) RegisterRest( c *rest.ServerContext ) {
  // Reference retrieval methods
  c.Handle( "/crs/{id}", r.CrsHandler ).Methods( "GET" )
  c.Handle( "/tiploc/{id}", r.TiplocHandler ).Methods( "GET" )

  c.Handle( "/toc", r.AllTocsHandler ).Methods( "GET" )
  c.Handle( "/toc/{id}", r.TocHandler ).Methods( "GET" )

  // Data import
  c.Handle( "/import", r.ImportHandler ).Methods( "POST" )
}
