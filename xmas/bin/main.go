package main

import (
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nre-feeds/xmas"
	"log"
)

func main() {
	err := kernel.Launch(&xmas.XmasService{})
	if err != nil {
		log.Fatal(err)
	}
}
