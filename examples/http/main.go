package main

import (
	"log"
	"time"

	"github.com/bh90210/healthz"
)

func main() {
	var healthCheck healthz.Check

	// setting struct values is optional.
	healthCheck.Liveness = "live"
	healthCheck.Readiness = "ready"
	healthCheck.Port = "8080"

	go func() {
		if err := healthCheck.Start(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		log.Println("ready")
		healthCheck.Ready()
		time.Sleep(time.Second * 60)
		log.Println("not ready")
		healthCheck.NotReady()
	}()

	if term := healthCheck.Terminating(); term == true {
		// do something
	}
}
