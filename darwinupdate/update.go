// Package that handles FTP updates
package darwinupdate

type DarwinUpdate struct {
  // The server name
  Server  string
  // The ftp user
  User    string
  // The ftp password for the NRE ftp server
  Pass    string
}

/*
func (u *DarwinUpdate) SetupSchedule( cr *cron.Cron, schedule string ) {
  cr.AddFunc( schedule, func () {
    if err := u.Update( true ); err != nil {
      log.Println( "Failed import:", err )
    }
  })
}

// Is an update required
func (u *DarwinUpdate) InitialImport() {
  if  (u.Ref != nil && importRequiredTimetable( u.Ref )) ||
      (u.TT != nil && importRequiredTimetable( u.TT )) {

    // Run in the background
    go func() {
      if err := u.Update( false ); err != nil {
        log.Println( "Failed import:", err )
      }
    }()

  }
}

// Update performs an update of all data
func (u *DarwinUpdate) Update( force bool ) error {
  return u.ftp( func( con *ftp.ServerConn ) error {

    if u.Ref != nil && (force || importRequiredTimetable( u.Ref ) ) {
      if err:= u.ReferenceUpdate( con ); err != nil {
        return err
      }
    }

    if u.TT != nil && (force || importRequiredTimetable( u.TT ) ) {
      if err:= u.TimetableUpdate( con ); err != nil {
        return err
      }
    }

    return nil
  })
}
*/
