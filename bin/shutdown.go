package main

import (
  "log"
  "os"
  "os/signal"
  "syscall"
)

func (c *Config) initShutdown() error {

  // Listen to signals & close the db before exiting
  // SIGINT for ^C, SIGTERM for docker stopping the container
  sigs := make( chan os.Signal, 1 )
  signal.Notify( sigs, syscall.SIGINT, syscall.SIGTERM )
  go func() {
    sig := <-sigs
    log.Println( "Signal", sig )

    c.cron.Stop()

    if( c.Database.reference != nil ) {
      c.Database.reference.Close()
    }

    if( c.Database.timetable != nil ) {
      c.Database.timetable.Close()
    }

    log.Println( "Database closed" )

    os.Exit( 0 )
  }()

  return nil
}
