package ldb

import (
  bolt "github.com/coreos/bbolt"
  "darwind3"
  "darwinref"
  "darwintimetable"
  "fmt"
  "github.com/peter-mount/golib/rest"
  "sort"
  "time"
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
  // Any reasons
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

    var services []*Service

    var messages []*darwind3.StationMessage

    now := time.Now()
    var nowt darwintimetable.WorkingTime
    nowt.Set( (now.Hour()*3600) + (now.Minute()*60) )
    next := now.Add( time.Hour )
    var hour darwintimetable.WorkingTime
    hour.Set( (next.Hour()*3600) + (next.Minute()*60) )

    if err := station.Update( func() error {
      // Station messages
      for _, id := range station.messages {
        if sm := d.Darwin.Messages.Get( id ); sm != nil {
          sm.Self = r.Self( fmt.Sprintf( "/live/message/%d", sm.ID ) )
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
        if len( services ) < 20 &&
           nowt.Compare( &s.Location.Times.Time ) &&
           s.Location.Times.Time.Compare( &hour ) {
          service := s.Clone()
          service.Self = r.Self( "/live/schedule/" + service.RID )
          services = append( services, service )
        }
      }
      return nil
    } ); err != nil {
      return err
    }

    // Resolve vias
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

    // Now resolve the Tiplocs
    res := &result{
      Services: services,
      Tiplocs: darwinref.NewLocationMap(),
      Tocs: darwinref.NewTocMap(),
      Messages: messages,
      Reasons: darwinref.NewReasonMap(),
      Date: now,
      Self: r.Self( "/ldb/boards/" + crs ),
    }

    if err := d.Reference.View( func( tx *bolt.Tx ) error {
      // Station details
      if sl, ok := d.Reference.GetCrs( tx, crs ); ok {
        for _, l := range sl {
          res.Station = append( res.Station, l.Tiploc )
          res.Tiplocs.AddTiploc( d.Reference, tx, l.Tiploc )
        }
      }

      // Tiplocs within the departures
      for _, s := range services {
        // Service & location tiplocs
        res.Tiplocs.AddTiploc( d.Reference, tx, s.Destination )
        res.Tiplocs.AddTiploc( d.Reference, tx, s.Location.Tiploc )
        // Toc running this service
        res.Tocs.AddToc( d.Reference, tx, s.Toc )

        // Tiploc in a cancel or late reason
        if s.CancelReason.Tiploc != "" {
          res.Tiplocs.AddTiploc( d.Reference, tx, s.CancelReason.Tiploc )
        }
        if s.LateReason.Tiploc != "" {
          res.Tiplocs.AddTiploc( d.Reference, tx, s.LateReason.Tiploc )
        }

        // Also add the reason
        if s.CancelReason.Reason > 0 {
          res.Reasons.Add( s.CancelReason.Reason, true, tx, d.Reference )
        }
        if s.LateReason.Reason > 0 {
          res.Reasons.Add( s.LateReason.Reason, false, tx, d.Reference )
        }
      }

      // Add any toc's from the locations in tiplocs
      res.Tocs.AddLocations( d.Reference, tx, res.Tiplocs )

      return nil
    }); err != nil {
      return err
    }

    // set Self on any maps
    res.Reasons.Self( r )
    res.Tiplocs.Self( r )
    res.Tocs.Self( r )

    // Return the response
    r.Status( 200 ).
      Value( res )
  }

  return nil
}
