// Package that handles FTP updates
package darwinupdate

import (
  "darwinref"
  "github.com/jlaffaye/ftp"
  "github.com/peter-mount/golib/rest"
  "gopkg.in/robfig/cron.v2"
  "log"
)

type DarwinUpdate struct {
  Ref    *darwinref.DarwinReference
  Pass    string
}

// RegisterRest registers the rest endpoints into a ServerContext
func (u *DarwinUpdate) Setup( c *rest.ServerContext, cr *cron.Cron ) {
  u.setupRest( c )
  u.setupCron( cr )
  log.Println( "FTP Client enabled" )

  u.initialImport()
}

func (u *DarwinUpdate) setupRest( c *rest.ServerContext ) {
  c.Handle( "/reference", u.ReferenceHandler ).Methods( "GET" )
}

func (u *DarwinUpdate) setupCron( cr *cron.Cron ) {
}

func (u *DarwinUpdate) initialImport() {
  importRef := u.Ref != nil && u.Ref.TimetableId() == ""

  if importRef {

    if err := u.ftp( func( con *ftp.ServerConn ) error {

      if importRef {
        if err:= u.ReferenceUpdate( con ); err != nil {
          return err
        }
      }

      return nil
    }); err != nil {
      log.Println( "Failed import:", err )
    }
    
  }
}
