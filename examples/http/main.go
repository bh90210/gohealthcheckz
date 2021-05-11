package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bh90210/healthz"
)

func init() {
	log.Println("init")
	go func() {
		log.Println("go")
		err := healthz.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	log.Println("finish init")
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	log.Println("ready")
	healthz.Ready()
	time.Sleep(time.Second * 60)
	log.Println("not ready")
	healthz.NotReady()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
