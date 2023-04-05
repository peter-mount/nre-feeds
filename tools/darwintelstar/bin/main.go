package main

import (
	"github.com/peter-mount/go-kernel/v2"
	telstar "github.com/peter-mount/nre-feeds/tools/darwintelstar"
	"log"
)

func main() {
	err := kernel.Launch(&telstar.Telstar{})
	if err != nil {
		log.Fatal(err)
	}
}
