package darwintimetable

import (
  "github.com/peter-mount/golib/rest"
  "log"
)

func (dt *DarwinTimetable) ImportHandler( r *rest.Rest ) error {
  // PruneSchedules first so we've removed anu old data before the import
  if _, err := dt.PruneSchedules(); err != nil {
    return err
  }

  log.Println( "DarwinTimetable import: started" )

  // Unmarshal the body, this actually does the import
  if err := r.Body( dt ); err != nil {
    return err
  }

  log.Println( "DarwinTimetable import: completed" )

  // PruneSchedules again to remove any old data we've just imported,
  // e.g. an old timetable file
  if _, err := dt.PruneSchedules(); err != nil {
    return err
  }

  r.Status( 200 ).
    Value( "ok" )
    return nil
}
