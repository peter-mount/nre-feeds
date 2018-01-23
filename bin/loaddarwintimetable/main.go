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
  var timetable darwintimetable.Timetable;
  err = xml.Unmarshal( data, &timetable )
  if err != nil {
    log.Fatal( err )
  }

  log.Println( "TimetableId", timetable.TimetableId )
  log.Println( "  Journey's", len( timetable.Journeys ) )

  log.Printf(  "\n%v\n", timetable.Journeys[ "201801048795721" ] )
  log.Printf(  "\n%v\n", timetable.Journeys[ "201801046763762" ] )

  log.Printf(  "\n%v\n", timetable.Journeys[ "201801047172127" ] )
  log.Printf(  "\n%v\n", timetable.Journeys[ "201801047174151" ] )

  log.Printf(  "\n%v\n", timetable.Journeys[ "201801048014393" ] )

  log.Printf(  "\n%v\n", timetable.Journeys[ "201801038010002" ] )

}
