// Load the darwin reference data file

package main

import (
  "darwinref"
  "encoding/xml"
  "flag"
  "io/ioutil"
  "log"
  "path/filepath"
)

func main() {
  log.Println( "loaddarwinref v0.1" )

  srcFile := flag.String( "f", "/ref.xml", "The config file to use" )

  flag.Parse()

  if srcFile == nil || *srcFile == "" {
    log.Fatal( "Source file -f required" )
  }

  filename, _ := filepath.Abs( *srcFile )
  log.Println( "Loading xml:", filename )

  data, err := ioutil.ReadFile( filename )
  if err != nil {
    log.Fatal( err )
  }

  log.Println( "Unmarshal" )
  var ref darwinref.DarwinReference;
  err = xml.Unmarshal( data, &ref )
  if err != nil {
    log.Fatal( err )
  }

  log.Printf( "Imported %v\n", &ref )

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
