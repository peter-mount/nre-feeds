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

type stationResult struct {
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
  // Cancellation or Late Reasons
  Reasons    *darwinref.ReasonMap     `json:"reasons"`
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
        // Ignore if it's departed
        include := !s.Location.Forecast.Departed

        if include {
          // Limit to max 20 departures and only if within the next hour
          include = len( services ) < 20
        }

        if include {
          include = nowt.Compare( &s.Location.Times.Time ) &&
            s.Location.Times.Time.Compare( &hour )
        }

        if include {
          service := s.Clone()
          // Point self to our proxy so we provide reference data as well
          service.Self = r.Self( "/service/" + service.RID )
          services = append( services, service )
        }
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

    res := &stationResult{
      Services: services,
      Tiplocs: darwinref.NewLocationMap(),
      Tocs: darwinref.NewTocMap(),
      Messages: messages,
      Reasons: darwinref.NewReasonMap(),
      Date: now,
      Self: r.Self( "/boards/" + crs ),
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

      if s.CancelReason.Reason > 0 {
        if reason, err := refClient.GetCancelledReason( s.CancelReason.Reason ); err != nil {
          return err
        } else if reason != nil {
          res.Reasons.AddReason( reason )
        }

        if s.CancelReason.Tiploc != "" {
          tiplocs[ s.CancelReason.Tiploc ] = nil
        }
      }

      if s.LateReason.Reason > 0 {
        if reason, err := refClient.GetLateReason( s.LateReason.Reason ); err != nil {
          return err
        } else if reason != nil {
          res.Reasons.AddReason( reason )
        }

        if s.LateReason.Tiploc != "" {
          tiplocs[ s.LateReason.Tiploc ] = nil
        }
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
