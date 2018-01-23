// Load the darwin reference data file

package main

import (
  "darwinref"
  "encoding/xml"
  "flag"
  "io/ioutil"
  "log"
  "os"
  "os/signal"
  "path/filepath"
  "syscall"
)

func main() {
  log.Println( "loaddarwinref v0.1" )

  dbFile := flag.String( "d", "/darwin.db", "The database file" )
  srcFile := flag.String( "f", "/ref.xml", "The config file to use" )

  flag.Parse()

  if srcFile == nil || *srcFile == "" {
    log.Fatal( "Source file -f required" )
  }

  var ref darwinref.DarwinReference;

  log.Printf( "Opening database %s\n", *dbFile )

  if err := ref.OpenDB( *dbFile ); err != nil {
    log.Fatal( err )
  }

  // Listen to signals & close the db before exiting
  // SIGINT for ^C, SIGTERM for docker stopping the container
  sigs := make( chan os.Signal, 1 )
  signal.Notify( sigs, syscall.SIGINT, syscall.SIGTERM )
  go func() {
    sig := <-sigs
    log.Println( "Signal", sig )

    ref.Close()
    log.Println( "Database closed" )

    os.Exit( 0 )
  }()

  log.Printf( "Importing %s\n", *srcFile )

  filename, _ := filepath.Abs( *srcFile )
  log.Println( "Loading xml:", filename )

  data, err := ioutil.ReadFile( filename )
  if err != nil {
    log.Fatal( err )
  }

  err = xml.Unmarshal( data, &ref )
  if err != nil {
    log.Fatal( err )
  }

  log.Println( "Reference file imported" )

  log.Printf( "Via %v\n", test( &ref, "AFK", "VICTRIE" ) )
  log.Printf( "Via %v\n", test( &ref, "AFK", "CANONST" ) )
  log.Printf( "Via %v\n", test( &ref, "AFK", "LNDNBDE" ) )
  log.Printf( "Via %v\n", test( &ref, "AFK", "CHRX" ) )
  log.Printf( "Via %v\n", test( &ref, "AFK", "BLFR" ) )
  log.Printf( "Via %v\n", test( &ref, "AFK", "MSTONEE" ) )

}

func test( ref *darwinref.DarwinReference, a string, d string ) []*darwinref.Via {
  vias, _ := ref.GetViaAt( a, d )
  return vias
}
