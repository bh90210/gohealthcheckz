package main

import (
	"log"
	"time"

	healthz "github.com/bh90210/healthz/grpc"
)

func main() {
	h, err := healthz.NewCheckGRPC(healthz.LivePath("live"))
	if err != nil {
		panic(err)
	}

	go func() {
		if err := h.Start(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		log.Println("not ready")
		h.NotReady()
		time.Sleep(time.Second * 5)
		log.Println("ready")
		h.Ready()
	}()

	if term := h.Terminating(); term == true {
		// do something
	}
}
