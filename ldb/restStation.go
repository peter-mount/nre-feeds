package ldb

import (
  "darwintimetable"
  "github.com/peter-mount/golib/rest"
  "sort"
  "time"
)

type result struct {
  Crs         string
  Services []*Service
  Date        time.Time
  Self        string
}

func (d *LDB) stationHandler( r *rest.Rest ) error {

  crs := r.Var( "crs" )

  station := d.GetStationCrs( crs )

  if station == nil {
    r.Status( 404 )
  } else {

    var services []*Service

    now := time.Now()
    var nowt darwintimetable.WorkingTime
    nowt.Set( (now.Hour()*3600) + (now.Minute()*60) )
    now = now.Add( time.Hour )
    var hour darwintimetable.WorkingTime
    hour.Set( (now.Hour()*3600) + (now.Minute()*60) )

    // Get the services from the station
    if err := station.Update( func() error {
      for _,s := range station.services {
        // Limit to max 20 departures and only if within the next hour
        if len( services ) < 20 &&
           nowt.Compare( &s.Location.Times.Time ) &&
           s.Location.Times.Time.Compare( &hour ) {
          service := s.Clone()
          service.Self = r.Self( "/ldb/service/" + service.RID )
          services = append( services, service )
        }
      }
      return nil
    } ); err != nil {
      return err
    }

    // sort into time order
    sort.SliceStable( services, func( i, j int ) bool {
      return services[ i ].Compare( services[ j ] )
    } )

    res := &result{
      Crs: crs,
      Services: services,
      Date: now,
      Self: r.Self( "/ldb/boards/" + crs ),
    }

    r.Status( 200 ).
      Value( res )
  }

  return nil
}
