// CIF Rest server
package main

import (
  "flag"
  "log"
)

func main() {
  log.Println( "darwin v0.1" )

  configFile := flag.String( "c", "", "The config file to use" )

  flag.Parse()

  if *configFile == "" {
    log.Fatal( "No default config defined, provide with -c" )
  }

  config := &Config{}

  if err := config.ReadFile( *configFile ); err != nil {
    log.Fatal( err )
  }

  if err := config.initCron(); err != nil {
    log.Fatal( err )
  }

  if err := config.initStats(); err != nil {
    log.Fatal( err )
  }

  if err := config.initServer(); err != nil {
    log.Fatal( err )
  }

  if err := config.initDb(); err != nil {
    log.Fatal( err )
  }

  if err := config.initFtp(); err != nil {
    log.Fatal( err )
  }

  if err := config.initShutdown(); err != nil {
    log.Fatal( err )
  }

  if err := config.start(); err != nil {
    log.Fatal( err )
  }
}
