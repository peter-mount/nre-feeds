package darwintimetable

import (
  "github.com/peter-mount/golib/rest"
  "log"
)

func (dt *DarwinTimetable) ImportHandler( r *rest.Rest ) error {
  log.Println( "DarwinTimetable import: started" )

  if err := r.Body( dt ); err != nil {
    return err
  }

  log.Println( "DarwinTimetable import: completed" )
  r.Status( 200 ).
    Value( "ok" )
    return nil
}
