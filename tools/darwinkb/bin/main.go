package main

import (
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nre-feeds/darwinkb/service"
	"log"
)

func main() {
	err := kernel.Launch(&service.DarwinKBService{})
	if err != nil {
		log.Fatal(err)
	}
}
