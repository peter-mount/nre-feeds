package main

import (
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nre-feeds/tools/darwintty"
	"log"
)

func main() {
	err := kernel.Launch(&darwintty.Server{})
	if err != nil {
		log.Fatal(err)
	}
}
