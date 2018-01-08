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
  log.Println( "LateRunningReasons", len( reference.LateRunningReasons ) )
  log.Println( "CancellationReasons", len( reference.CancellationReasons ) )
  log.Println( "CISSource's", len( reference.CISSource ) )

  log.Println( "---------")

  var ref *darwinref.DarwinReference = reference.Decode()
  log.Println( "             Tiploc", len( ref.Tiploc ) )
  log.Println( "                Crs", len( ref.Crs ) )
  log.Println( "              Toc's", len( ref.Toc ) )
  log.Println( " LateRunningReasons", len( ref.LateRunningReasons ) )
  log.Println( "CancellationReasons", len( ref.CancellationReasons ) )
  log.Println( "          CISSource", len( ref.CISSource ) )
}
