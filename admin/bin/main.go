package main

import (
  "github.com/peter-mount/go-kernel"
  "github.com/peter-mount/nre-feeds/admin/messages"
  "log"
)

func main() {
  err := kernel.Launch(
    &messages.Messages{},
  )
  if err != nil {
    log.Fatal(err)
  }
}
