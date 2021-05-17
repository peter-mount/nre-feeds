package main

import (
	"github.com/peter-mount/golib/kernel"
	"github.com/peter-mount/nre-feeds/darwingraph/service"
	"log"
)

func main() {
	err := kernel.Launch(
		&service.DarwinGraphService{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
