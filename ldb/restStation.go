package ldb

import (
  "darwind3"
  "darwinref"
//  "fmt"
  "github.com/peter-mount/golib/rest"
  "sort"
  "time"
  "util"
)

type result struct {
  // The departures
  Services []*Service                 `json:"departures"`
  // Details about this station
  Station  []string                   `json:"station"`
  // Map of Tiploc's
  Tiplocs    *darwinref.LocationMap   `json:"tiploc"`
  // Map of Toc's
  Tocs       *darwinref.TocMap        `json:"toc"`
  // StationMessages
  Messages []*darwind3.StationMessage `json:"messages"`
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

    d3Client := &darwind3.DarwinD3Client{ Url: d.Darwin }
    refClient := &darwinref.DarwinRefClient{ Url: d.Reference }

    var services []*Service

    var messages []*darwind3.StationMessage

    now := time.Now()
    var nowt util.WorkingTime
    nowt.Set( (now.Hour()*3600) + (now.Minute()*60) )
    next := now.Add( time.Hour )
    var hour util.WorkingTime
    hour.Set( (next.Hour()*3600) + (next.Minute()*60) )

    if err := station.Update( func() error {
      // Station messages
      for _, id := range station.messages {
        if sm, _ := d3Client.GetStationMessage( id ); sm != nil {
          messages = append( messages, sm )
        }
      }

      // Get the services from the station
      var sa []*Service
      for _,s := range station.services {
        sa = append( sa, s )
      }

      // sort into time order
      sort.SliceStable( sa, func( i, j int ) bool {
        return sa[ i ].Compare( sa[ j ] )
      } )

      for _, s := range sa {
        // Limit to max 20 departures and only if within the next hour
        //if len( services ) < 20 &&
        //   nowt.Compare( &s.Location.Times.Time ) &&
        //   s.Location.Times.Time.Compare( &hour ) {
          service := s.Clone()
          service.Self = r.Self( "/live/schedule/" + service.RID )
          services = append( services, service )
        //}
      }
      return nil
    } ); err != nil {
      return err
    }

    // Resolve vias
    /*
    for _, s := range services {
      sv := d.Darwin.GetSchedule( s.RID )
      if sv != nil {sv.View( func() error {
        // Find our Location
        found := false
        var locs []string
        for _, l := range sv.Locations {
          if found {
            locs = append( locs, l.Tiploc )
          } else if l.Equals( s.Location ) {
            found = true
          }
        }

        if len( locs ) > 0 {
          via := d.Reference.ResolveVia( crs, s.Destination, locs )
          if via != nil {
            s.Via = via.Text
          }
        }
        return nil
      })
      }
    }
    */

    res := &result{
      Services: services,
      Tiplocs: darwinref.NewLocationMap(),
      Tocs: darwinref.NewTocMap(),
      Messages: messages,
      Date: now,
      Self: r.Self( "/ldb/boards/" + crs ),
    }

    // Set of tiplocs
    tiplocs := make( map[string]interface{} )

    // Station details
    if sl, _ := refClient.GetCrs( crs ); sl != nil {
      for _, l := range sl.Tiploc {
        res.Station = append( res.Station, l.Tiploc )
        tiplocs[ l.Tiploc ] = nil
      }
    }

    // Tiplocs within the departures
    for _, s := range services {

      // Destination & location tiplocs
      tiplocs[ s.Destination ] = nil
      tiplocs[ s.Location.Tiploc ] = nil

      // Toc running this service
      refClient.AddToc( res.Tocs, s.Toc )

      // Tiploc in a cancel or late reason
      if s.CancelReason.Tiploc != "" {
        tiplocs[ s.CancelReason.Tiploc ] = nil
      }

      if s.LateReason.Tiploc != "" {
        tiplocs[ s.LateReason.Tiploc ] = nil
      }
    }

    // Now resolve the tiplocs en-masse and resolve the toc's at the same time
    if locs, _ := refClient.GetTiplocsMapKeys( tiplocs ); locs != nil {
      res.Tiplocs.AddAll( locs )

      for _, l := range locs {
        refClient.AddToc( res.Tocs, l.Toc )
      }
    }

    r.Status( 200 ).
      Value( res )
  }

  return nil
}
