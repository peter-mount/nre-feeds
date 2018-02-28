package ldb

import (
  "darwind3"
  "darwinref"
  "github.com/peter-mount/golib/rest"
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

    // We want everything for the next hour
    now := time.Now()
    from := util.WorkingTime_FromTime( now )
    to := util.WorkingTime_FromTime( now.Add( time.Hour ) )

    services := station.GetServices( from, to )

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
      Messages: station.GetMessages( d3Client ),
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

      // Cancellation reason
      if s.CancelReason.Reason > 0 {
        if reason, _ := refClient.GetCancelledReason( s.CancelReason.Reason ); reason != nil {
          res.Reasons.AddReason( reason )
        }

        if s.CancelReason.Tiploc != "" {
          tiplocs[ s.CancelReason.Tiploc ] = nil
        }
      }

      // Late reason
      if s.LateReason.Reason > 0 {
        if reason, _ := refClient.GetLateReason( s.LateReason.Reason ); reason != nil {
          res.Reasons.AddReason( reason )
        }

        if s.LateReason.Tiploc != "" {
          tiplocs[ s.LateReason.Tiploc ] = nil
        }
      }

      // Set self to point to our service endpoint
      s.Self = r.Self( "/service/" + s.RID )
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
