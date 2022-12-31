package main

import (
	"github.com/peter-mount/go-kernel/v2"
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
