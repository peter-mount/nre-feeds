// Package that handles FTP updates
package darwinupdate

import (
  "darwinref"
  "darwintimetable"
  "github.com/jlaffaye/ftp"
  "github.com/peter-mount/golib/rest"
  "gopkg.in/robfig/cron.v2"
  "log"
  "time"
)

type DarwinUpdate struct {
  // DarwinReference instance or nil
  Ref    *darwinref.DarwinReference
  // DarwinTimetable instance or nil
  TT     *darwintimetable.DarwinTimetable
  // The server name
  Server  string
  // The ftp user
  User    string
  // The ftp password for the NRE ftp server
  Pass    string
}

func (u *DarwinUpdate) SetupRest( c *rest.ServerContext ) {
  c.Handle( "/reference", u.ReferenceHandler ).Methods( "GET" )
  c.Handle( "/timetable", u.TimetableHandler ).Methods( "GET" )

  // Expose Update()
  c.Handle( "/all", func( r *rest.Rest ) error {
    if err := u.Update( true ); err != nil {
      return err
    }

    r.Status( 200 ).Value( "ok" )
    return nil
  })
}

func (u *DarwinUpdate) SetupSchedule( cr *cron.Cron, schedule string ) {
  cr.AddFunc( schedule, func () {
    if err := u.Update( true ); err != nil {
      log.Println( "Failed import:", err )
    }
  })
}

// Is an update required
func importRequiredTimetable( v interface{ TimetableId() string } ) bool {
  // Import if no TimetableId
  if v.TimetableId() == "" {
    return true
  }

  // Import if TimetableId is older than the current day
  limit := time.Now().Truncate( 24 * time.Hour )
  tid, err := time.Parse( "20060102150405", v.TimetableId() )
  // Error then force import as tid is invalid
  return err != nil || tid.Before( limit )
}

func (u *DarwinUpdate) InitialImport() {
  if  (u.Ref != nil && importRequiredTimetable( u.Ref )) ||
      (u.TT != nil && importRequiredTimetable( u.TT )) {
    if err := u.Update( false ); err != nil {
      log.Println( "Failed import:", err )
    }
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
