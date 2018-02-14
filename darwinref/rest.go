package darwinref

import (
  "github.com/peter-mount/golib/rest"
)

// RegisterRest registers the rest endpoints into a ServerContext
func (r DarwinReference) RegisterRest( c *rest.ServerContext ) {

  c.Handle( "/reason/cancelled", r.AllReasonCancelHandler ).Methods( "GET" )
  c.Handle( "/reason/cancelled/{id}", r.ReasonCancelHandler ).Methods( "GET" )

  c.Handle( "/reason/late", r.AllReasonLateHandler ).Methods( "GET" )
  c.Handle( "/reason/late/{id}", r.ReasonLateHandler ).Methods( "GET" )

  c.Handle( "/via/{at}/{dest}/{loc1}", r.ViaHandler ).Methods( "GET" )
  c.Handle( "/via/{at}/{dest}/{loc1}/{loc2}", r.ViaHandler ).Methods( "GET" )

  // Reference retrieval methods
  c.Handle( "/crs/{id}", r.CrsHandler ).Methods( "GET" )
  c.Handle( "/crs", r.AllCrsHandler ).Methods( "GET" )

  c.Handle( "/tiploc", r.TiplocsHandler ).Methods( "POST" )
  c.Handle( "/tiploc/{id}", r.TiplocHandler ).Methods( "GET" )

  c.Handle( "/toc", r.AllTocsHandler ).Methods( "GET" )
  c.Handle( "/toc/{id}", r.TocHandler ).Methods( "GET" )

  // Data import
  c.Handle( "/import", r.ImportHandler ).Methods( "POST" )
}
