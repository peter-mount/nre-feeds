package darwinref

import (
  "github.com/peter-mount/golib/rest"
)

// An entry in the request object
type viaResolveRequest struct {
  // CRS of the location we want to show a via
  Crs           string    `json:"crs"`
  // Destination tiploc
  Destination   string    `json:"destination"`
  // Tiplocs of journey after this location to search
  Tiplocs     []string    `json:"tpls"`
}

// viaResolveHandler resolves The via(s) for a set of schedules
func (dr *DarwinReference) viaResolveHandler( r *rest.Rest ) error {

  // The query
  queries := make( map[string]*viaResolveRequest )

  // The response
  response := make( map[string]*Via )

  // Run the queries
  if err := r.Body( &queries ); err != nil {

    // Fail safe by returning 500 but still a {} object
    r.Status( 500 ).Value( response )

  } else {

    for rid, query := range queries {
      if via := dr.ResolveVia( query.Crs, query.Destination, query.Tiplocs ); via != nil {
        via.SetSelf( r )
        response[ rid ] = via
      }
    }

    r.Status( 200 ).Value( response )
  }

  return nil
}
