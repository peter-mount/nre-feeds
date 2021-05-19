package main

import (
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/nre-feeds/darwingraph/maps"
	"github.com/peter-mount/nre-feeds/darwingraph/service"
	"log"
)

func main() {
	err := kernel.Launch(
		&service.DarwinGraphService{},
		&maps.UKMap{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
