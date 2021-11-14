package main

import (
	"log"
	"time"

	"github.com/bh90210/healthz"
)

func main() {
	h := healthz.NewCheck(healthz.OptionsLivePath("live"),
		healthz.OptionsReadyPath("ready"), healthz.OptionsPort("8080"))

	go func() {
		if err := h.Start(); err != nil {
			panic(err)
		}
	}()

	go func() {
		log.Println("not ready")
		h.NotReady()
		time.Sleep(time.Second * 5)
		log.Println("ready")
		h.Ready()
	}()

	if h.Terminating() {
		// do something
	}
}
