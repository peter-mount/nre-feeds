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
  var reference darwinref.PportTimetableRef;
  err = xml.Unmarshal( data, &reference )
  if err != nil {
    log.Fatal( err )
  }

  log.Println( "Locations", len( reference.Locations ) )
  log.Println( "Toc's", len( reference.Toc ) )
  log.Println( "LateRunningReasons", len( reference.LateRunningReasons.Reason ) )
  log.Println( "CancellationReasons", len( reference.CancellationReasons.Reason ) )
  log.Println( "CISSource's", len( reference.CISSource ) )
}
