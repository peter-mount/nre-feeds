// CIF Rest server
package main

import (
  bolt "github.com/coreos/bbolt"
  "darwinref"
  "flag"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"
)

type server struct {
  db       *bolt.DB
  ref       darwinref.DarwinReference
  server   *rest.Server
}

func main() {
  log.Println( "darwin v0.1" )

  // TODO use this to protect /importCIF endpoint
  //writeSecret := flag.String( "s", "", "The write secret" )

  dbFile := flag.String( "d", "/database.db", "The database file" )

  // Port for the webserver
  port := flag.Int( "p", 8080, "Port to use" )

  flag.Parse()

  stats := statistics.Statistics{ Log: true }
  stats.Configure()

  server := &server{}

  if err := server.openDB( *dbFile ); err != nil {
    log.Fatal( err )
  }

  if err := server.ref.UseDB( server.db ); err != nil {
    log.Fatal( err )
  }

  // Listen to signals & close the db before exiting
  // SIGINT for ^C, SIGTERM for docker stopping the container
  sigs := make( chan os.Signal, 1 )
  signal.Notify( sigs, syscall.SIGINT, syscall.SIGTERM )
  go func() {
    sig := <-sigs
    log.Println( "Signal", sig )

    server.ref.Close()
    server.db.Close()
    log.Println( "Database closed" )

    os.Exit( 0 )
  }()

  server.server = rest.NewServer( *port )

  server.ref.RegisterRest( server.server.Context( "/ref" ) )

  server.server.Start()
}

func (s *server) openDB( dbFile string ) error {
  if db, err := bolt.Open( dbFile, 0666, &bolt.Options{
    Timeout: 5 * time.Second,
    } ); err != nil {
      return err
  } else {
    s.db = db
    return nil
  }
}
