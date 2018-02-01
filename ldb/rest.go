package ldb

import (
  "github.com/peter-mount/golib/rest"
)

func (d *LDB) RegisterRest( c *rest.ServerContext ) {
  c.Handle( "/boards/{crs}", d.stationHandler ).Methods( "GET" )
}
