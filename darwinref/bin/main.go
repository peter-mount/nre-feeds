package main

import (
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/nre-feeds/darwinref/service"
	"github.com/peter-mount/nre-feeds/darwinref/update"
	"log"
)

func main() {
	err := kernel.Launch(
		&service.DarwinRefService{},
		&update.ReferenceUpdateService{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
