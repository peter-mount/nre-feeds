package ldb

import (
  "github.com/peter-mount/golib/rest"
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

    // Copy the service slice inside the loc
    if err := station.Update( func() error {
      services = append( []*Service(nil), station.services... )
      return nil
    } ); err != nil {
      return err
    }

    res := &result{ Crs: crs, Services: services }

    r.Status( 200 ).
      Value( res )
  }

  return nil
}
