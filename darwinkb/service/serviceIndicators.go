package service

import (
  "github.com/peter-mount/golib/rest"
  "log"
)

func (d *DarwinKBService) GetServiceIndicatorsHandler( r *rest.Rest ) error {
  data, err := d.darwinkb.GetServiceIndicators()
  if err != nil {
    log.Println( err )
    return err
  }

  if data == nil {
    r.Status( 404 )
  } else {
    r.Status( 200 ).
      JSON().
      Writer().
      Write( data )
  }

  return nil
}

func (d *DarwinKBService) GetServiceIndicatorHandler( r *rest.Rest ) error {
  data, err := d.darwinkb.GetServiceIndicator( r.Var( "id" ) )
  if err != nil {
    log.Println( err )
    return err
  }

  if data == nil {
    r.Status( 404 )
  } else {
    r.Status( 200 ).
      JSON().
      Writer().
      Write( data )
  }

  return nil
}
