package update

import (
  "github.com/peter-mount/nre-feeds/darwind3"
  "log"
)

// timetableUpdateListener listens for real time updates for when new reference
// data is made available.
func (d *TimetableUpdateService) timetableUpdateListener( e *darwind3.DarwinEvent ) {

  if e.TimeTableId != nil && e.TimeTableId.TTFile != "" {
    log.Printf("New timetable %s", e.TimeTableId.TimeTableId )

    err := d.updateTimetable( e.TimeTableId )
    if err != nil {
      log.Printf( "Failed to import %s: %v", e.TimeTableId.TimeTableId, err )
    } else {
      log.Printf( "Imported %s", e.TimeTableId.TimeTableId )
    }
  }

}
