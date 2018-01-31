package ldb

import (
  bolt "github.com/coreos/bbolt"
  "darwinref"
  "darwintimetable"
  "github.com/peter-mount/golib/rest"
  "sort"
  "time"
)

type result struct {
  // The CRS of this station
  Crs         string                  `json:"crs"`
  // The departures
  Services []*Service                 `json:"departures"`
  // Map of Tiploc's
  Tiplocs    *darwinref.LocationMap   `json:"tiploc"`
  // Map of Toc's
  Tocs       *darwinref.TocMap        `json:"toc"`
  // The date of this request
  Date        time.Time               `json:"date"`
  // The URL of this departure board
  Self        string                  `json:"self"`
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
    next := now.Add( time.Hour )
    var hour darwintimetable.WorkingTime
    hour.Set( (next.Hour()*3600) + (next.Minute()*60) )

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

    // Now resolve the Tiplocs
    tiplocs := darwinref.NewLocationMap()
    tocs := darwinref.NewTocMap()
    if err := d.Reference.View( func( tx *bolt.Tx ) error {
      for _, s := range services {
        // Service & location tiplocs
        tiplocs.AddTiploc( d.Reference, tx, s.Destination )
        tiplocs.AddTiploc( d.Reference, tx, s.Location.Tiploc )
        // Toc running this service
        tocs.AddToc( d.Reference, tx, s.Toc )
      }

      // Add any toc's from the locations in tiplocs
      tocs.AddLocations( d.Reference, tx, tiplocs )

      return nil
    }); err != nil {
      return err
    }

    tiplocs.Self( r )
    tocs.Self( r )
    
    res := &result{
      Crs: crs,
      Services: services,
      Tiplocs: tiplocs,
      Tocs: tocs,
      Date: now,
      Self: r.Self( "/ldb/boards/" + crs ),
    }

    r.Status( 200 ).
      Value( res )
  }

  return nil
}
