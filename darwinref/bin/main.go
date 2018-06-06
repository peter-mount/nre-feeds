package main

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/nre-feeds/darwinref/service"
  "log"
)

func main() {
  err := kernel.Launch( &service.DarwinRefService{} )
  if err != nil {
    log.Fatal( err )
  }
}
