package main

import (
	"github.com/bazookajoe1/metrics-collector/internal/agent"
	"log"
)

func main() {
	log.Fatal(agent.RunAgent())
}
