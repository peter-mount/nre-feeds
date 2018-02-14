package bin

import (
  "log"
  "os"
  "os/signal"
  "syscall"
)

// initShutdown adds signal handlers to allow clean shutdown within a Docker container
func (c *Config) initShutdown( close func() ) error {

  // Listen to signals & close the db before exiting
  // SIGINT for ^C, SIGTERM for docker stopping the container
  sigs := make( chan os.Signal, 1 )
  signal.Notify( sigs, syscall.SIGINT, syscall.SIGTERM )
  go func() {
    sig := <-sigs
    log.Println( "Signal", sig )

    c.Cron.Stop()

    if close != nil {
      close()
    }

    log.Println( "Database closed" )

    os.Exit( 0 )
  }()

  return nil
}
