package ldb

import (
  "darwind3"
  "darwinref"
  "github.com/peter-mount/golib/rest"
  "time"
)

type serviceResult struct {
  // The service RID
  RID           string                `json:"rid"`
  // Origin
  Origin       *darwind3.Location     `json:"origin"`
  // Destination
  Destination  *darwind3.Location     `json:"destination"`
  // The service
  Service      *darwind3.Schedule     `json:"service"`
  // Map of Tiploc's
  Tiplocs      *darwinref.LocationMap `json:"tiploc"`
  // Map of Toc's
  Tocs         *darwinref.TocMap      `json:"toc"`
  // The date of this request
  Date          time.Time             `json:"date"`
  // The URL of this departure board
  Self          string                `json:"self"`
}

// serviceHandler proxies the service from d3 but fills in the required
// details of tiplocs, toc etc
func (d *LDB) serviceHandler( r *rest.Rest ) error {

  rid := r.Var( "rid" )

  d3Client := &darwind3.DarwinD3Client{ Url: d.Darwin }

  if service, err := d3Client.GetSchedule( rid ); err != nil {
    return err
  } else if service == nil {
    r.Status( 404 )
  } else {

    refClient := &darwinref.DarwinRefClient{ Url: d.Reference }

    res := &serviceResult{
      RID: rid,
      Service: service,
      Tiplocs: darwinref.NewLocationMap(),
      Tocs: darwinref.NewTocMap(),
      Date: time.Now(),
      Self: r.Self( "/service/" + rid ),
    }

    // The origin & destination are the first & last locations in the schedule
    if len( service.Locations ) > 0 {
      res.Origin = service.Locations[ 0 ]
      res.Destination = service.Locations[ len( service.Locations ) - 1 ]
    }

    // Set of tiplocs
    tiplocs := make( map[string]interface{} )
    for _, l := range service.Locations {
      tiplocs[ l.Tiploc ] = nil
    }

    // Toc running this service
    refClient.AddToc( res.Tocs, service.Toc )

    // Tiploc in a cancel or late reason
    if service.CancelReason.Tiploc != "" {
      tiplocs[ service.CancelReason.Tiploc ] = nil
    }

    if service.LateReason.Tiploc != "" {
      tiplocs[ service.LateReason.Tiploc ] = nil
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