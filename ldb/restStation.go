package ldb

import (
  "github.com/peter-mount/golib/rest"
  "sort"
)

type result struct {
  Crs     string
  Services    []*Service
}

func (d *LDB) stationHandler( r *rest.Rest ) error {

  crs := r.Var( "crs" )

  station := d.GetStationCrs( crs )

  if station == nil {
    r.Status( 404 )
  } else {

    var services []*Service

    // Get the services from the station
    if err := station.Update( func() error {
      for _,s := range station.services {
        services = append( services, s )
      }
      return nil
    } ); err != nil {
      return err
    }

    // sort into time order
    sort.SliceStable( services, func( i, j int ) bool {
      return services[ i ].Compare( services[ j ] )
    } )

    res := &result{ Crs: crs, Services: services }

    r.Status( 200 ).
      Value( res )
  }

  return nil
}
