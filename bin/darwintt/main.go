// darwintt Microservice
package main

import (
  "bin"
  "darwintimetable"
  "flag"
  "log"
)

func main() {
  log.Println( "darwintt v0.1" )

  configFile := flag.String( "c", "", "The config file to use" )

  flag.Parse()

  if *configFile == "" {
    log.Fatal( "No default config defined, provide with -c" )
  }

  config := &bin.Config{}

  if err := config.ReadFile( *configFile ); err != nil {
    log.Fatal( err )
  }

  if err := config.InitCron(); err != nil {
    log.Fatal( err )
  }

  if err := config.InitServer(); err != nil {
    log.Fatal( err )
  }

  if err := config.InitStats(); err != nil {
    log.Fatal( err )
  }

  if err := config.InitDb(); err != nil {
    log.Fatal( err )
  }

  // Reference database
  tt := &darwintimetable.DarwinTimetable{}
  config.DbPath( &config.Database.Timetable, "dwtt.db" )
  if err := tt.OpenDB( config.Database.Timetable ); err != nil {
    log.Fatal( err )
  }
  tt.RegisterRest( config.Server.Ctx )

  // Enable ftp to auto update
  if err := config.InitFtp(); err != nil {
    log.Fatal( err )
  } else if config.Ftp.Enabled {
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

  if err := config.InitShutdown( tt.Close ); err != nil {
    log.Fatal( err )
  }

  if err := config.Start(); err != nil {
    log.Fatal( err )
  }
}
