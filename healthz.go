// Package healthz is a small & simple to use library for liveness & readiness Kubernetes checks (gRPC included).
package healthz

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type state int

const (
	notready state = iota
	ready
)

var s state

// Start .
func Start() error {
	http.Handle("/ready", readiness())
	http.Handle("/live", liveness())
	return http.ListenAndServe(":6080", nil)
}

// Ready .
func Ready() {
	s = ready
}

// NotReady .
func NotReady() {
	s = notready
}

// Terminating .
func Terminating() bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()
	return <-done
}

func liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s == ready {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func readiness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s == ready {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
