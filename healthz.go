// Package healthz is a tiny & simple to use library for liveness & readiness Kubernetes checks (gRPC included).
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

// Start starts the healthcheck http server. It should be called at the start of your application.
// It is a blocking function.
func Start() error {
	http.Handle("/ready", readiness())
	http.Handle("/live", liveness())
	return http.ListenAndServe(":6080", nil)
}

// Ready sets the state of service to ready. State's default value is false.
// You have to manually enabled whenever app is ready to service requests.
func Ready() {
	s = ready
}

// NotReady sets the state to notready.
func NotReady() {
	s = notready
}

// Terminating starts a go routine waiting for SIGINT & SIGTERM signals.
// Returns true when Kubernetes sends a termination signal to the pod.
// It is a blocking function.
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
