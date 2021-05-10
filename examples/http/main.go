package main

import (
	"log"

	"github.com/bh90210/healthz"
)

func main() {
	log.Println("test starting")
	req := make(chan healthz.State)
	res := make(chan bool)
	healthz.LivenessReadiness(req, res, func() {})

	log.Println("block")
	block := make(chan bool)
	<-block
}
