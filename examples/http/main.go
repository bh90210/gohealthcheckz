package main

import (
	"log"
	"time"

	"github.com/bh90210/healthz"
)

func init() {
	go func() {
		if err := healthz.Start(); err != nil {
			log.Fatalln(err)
		}
	}()
}

func main() {
	go func() {
		log.Println("ready")
		healthz.Ready()
		time.Sleep(time.Second * 60)
		log.Println("not ready")
		healthz.NotReady()
	}()

	if term := healthz.Terminating(); term == true {
		// do something
	}
}
