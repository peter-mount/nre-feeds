package main

import (
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nre-feeds/tools/darwintelstar/boards"
	"github.com/peter-mount/nre-feeds/tools/darwintelstar/index"
	"log"
)

func main() {
	err := kernel.Launch(
		&boards.Departures{},
		&index.Index{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
