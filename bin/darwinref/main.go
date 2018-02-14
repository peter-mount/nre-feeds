// CIF Rest server
package main

import (
  "bin"
  "darwinref"
  "flag"
  "log"
)

func main() {
  log.Println( "darwinref v0.1" )

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
  reference := &darwinref.DarwinReference{}
  config.DbPath( &config.Database.Reference, "dwref.db" )
  if err := reference.OpenDB( config.Database.Reference ); err != nil {
    log.Fatal( err )
  }
  reference.RegisterRest( config.Server.Ctx )

  // Enable ftp to auto update
  if err := config.InitFtp(); err != nil {
    log.Fatal( err )
  } else if config.Ftp.Enabled {
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

  if err := config.InitShutdown( reference.Close ); err != nil {
    log.Fatal( err )
  }

  if err := config.Start(); err != nil {
    log.Fatal( err )
  }
}
