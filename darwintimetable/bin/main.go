package main

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/nre-feeds/darwintimetable/service"
  "log"
)

func main() {
  err := kernel.Launch( &service.DarwinTimetableService{} )
  if err != nil {
    log.Fatal( err )
  }
}
