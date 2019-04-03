package main

import (
	"github.com/peter-mount/golib/kernel"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwind3/service"
	"log"
)

func main() {
	err := kernel.Launch(
		&bin.Graphite{},
		&service.DarwinD3Service{})
	if err != nil {
		log.Fatal(err)
	}
}
