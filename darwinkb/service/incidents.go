package service

import (
  "github.com/peter-mount/golib/rest"
  "log"
)

func (d *DarwinKBService) IncidentsHandler( r *rest.Rest ) error {
  data, err := d.darwinkb.GetIncidents()
  if err != nil {
    log.Println( err )
    return err
  }

  if data == nil {
    r.Status( 404 )
  } else {
    r.Status( 200 ).
    Writer().
    Write( data )
  }

  return nil
}

func (d *DarwinKBService) IncidentHandler( r *rest.Rest ) error {
  data, err := d.darwinkb.GetIncident( r.Var( "id" ) )
  if err != nil {
    log.Println( err )
    return err
  }

  if data == nil {
    r.Status( 404 )
  } else {
    r.Status( 200 ).
    Writer().
    Write( data )
  }

  return nil
}
