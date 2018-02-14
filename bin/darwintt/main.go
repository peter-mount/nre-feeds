// darwintt Microservice
package main

import (
  "bin"
  "darwintimetable"
  "log"
)

func main() {
  log.Println( "darwintt v0.1" )

  bin.RunApplication( app )
}

func app( config *bin.Config ) ( func(), error ) {

  // Reference database
  tt := &darwintimetable.DarwinTimetable{}

  config.DbPath( &config.Database.Timetable, "dwtt.db" )

  if err := tt.OpenDB( config.Database.Timetable ); err != nil {
    return nil, err
  }

  tt.RegisterRest( config.Server.Ctx )

  // Enable ftp to auto update
  if err := config.InitFtp(); err != nil {
    return nil, err
  }

  if config.Ftp.Enabled {
    // Scheduled updates
    if config.Ftp.Schedule != "" {
      config.Cron.AddFunc( config.Ftp.Schedule, func () {
        if err := config.Ftp.Update.TimetableUpdate( tt ); err != nil {
          log.Println( "Failed import:", err )
        }
      })
      log.Println( "Auto Update using:", config.Ftp.Schedule )
    }

    // Initial import required?
    if config.Ftp.Update.ImportRequiredTimetable( tt ) {
      config.Ftp.Update.TimetableUpdate( tt )
    }
  }

  return tt.Close, nil
}
