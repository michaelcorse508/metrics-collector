package main

import (
	"github.com/bazookajoe1/metrics-collector/internal/server"
	"log"
)

func main() {
	log.Fatal(server.RunServer())
}
