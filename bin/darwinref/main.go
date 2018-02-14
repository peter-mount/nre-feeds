// darwinref Microservice
package main

import (
  "bin"
  "darwinref"
  "log"
)

func main() {
  log.Println( "darwinref v0.1" )

  bin.RunApplication( app )
}

func app( config *bin.Config ) ( func(), error ) {

  // Reference database
  reference := &darwinref.DarwinReference{}
  config.DbPath( &config.Database.Reference, "dwref.db" )
  if err := reference.OpenDB( config.Database.Reference ); err != nil {
    return nil, err
  }

  reference.RegisterRest( config.Server.Ctx )

  // Enable ftp to auto update
  if err := config.InitFtp(); err != nil {
    return nil, err
  }

  if config.Ftp.Enabled {
    // Scheduled updates
    if config.Ftp.Schedule != "" {
      config.Cron.AddFunc( config.Ftp.Schedule, func () {
        if err := config.Ftp.Update.ReferenceUpdate( reference ); err != nil {
          log.Println( "Failed import:", err )
        }
      })
      log.Println( "Auto Update using:", config.Ftp.Schedule )
    }

    // Initial import required?
    if config.Ftp.Update.ImportRequiredTimetable( reference ) {
      config.Ftp.Update.ReferenceUpdate( reference )
    }
  }

  return reference.Close, nil
}
