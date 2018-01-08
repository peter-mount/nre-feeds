// Load the darwin reference data file

package main

import (
  "darwintimetable"
  "encoding/xml"
  "flag"
  "io/ioutil"
  "log"
  "path/filepath"
)

func main() {
  log.Println( "loaddarwintimetable v0.1" )

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
  var timetable darwintimetable.PportTimetable;
  err = xml.Unmarshal( data, &timetable )
  if err != nil {
    log.Fatal( err )
  }

  log.Println( "TimetableId", timetable.TimetableId )
  log.Println( "  Journey's", len( timetable.Journeys ) )

  for i := 0; i<10; i++ {
    log.Printf(  "Journey %d\n%v\n", i, timetable.Journeys[ i ] )
  }
}
