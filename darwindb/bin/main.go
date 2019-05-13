package main

import (
	"github.com/peter-mount/golib/kernel"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/darwindb/service"
	"log"
)

func main() {
	err := kernel.Launch(
		&bin.Graphite{},
		&service.DarwinDBService{})
	if err != nil {
		log.Fatal(err)
	}
}
