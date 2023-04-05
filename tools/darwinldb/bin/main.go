package main

import (
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nre-feeds/bin"
	"github.com/peter-mount/nre-feeds/ldb/service"
	"log"
)

func main() {
	err := kernel.Launch(
		&bin.Graphite{},
		&service.LDBService{})
	if err != nil {
		log.Fatal(err)
	}
}
