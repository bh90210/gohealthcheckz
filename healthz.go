// Package healthz is a small & simple to use library for liveness & readiness Kubernetes checks (gRPC included).
package healthz

import (
	"net/http"
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

func Ready() {
	s = ready
}

func NotReady() {
	s = notready
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

func liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s == ready {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

// startup

func terminating() {

}
