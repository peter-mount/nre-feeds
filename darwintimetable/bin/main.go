package main

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/nre-feeds/darwintimetable/service"
  "github.com/peter-mount/nre-feeds/darwintimetable/update"
  "log"
)

func main() {
  err := kernel.Launch(
    &service.DarwinTimetableService{},
    &update.TimetableUpdateService{},
  )
  if err != nil {
    log.Fatal( err )
  }
}
