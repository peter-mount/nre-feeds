package main

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/nre-feeds/darwinkb/service"
  "log"
)

func main() {
  err := kernel.Launch( &service.DarwinKBService{} )
  if err != nil {
    log.Fatal( err )
  }
}
