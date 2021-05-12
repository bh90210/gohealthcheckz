// Package healthz is a tiny & simple to use library for liveness & readiness Kubernetes checks (gRPC included).
package healthz

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type state int

const (
	notready state = iota
	ready
)

type Check struct {
	Liveness  string
	Readiness string
	Port      string
	state     state
}

// Start starts the healthcheck http server. It should be called at the start of your application.
// It is a blocking function.
func (h *Check) Start() error {
	if len(h.Liveness) == 0 {
		h.Liveness = "/live"
	} else if !strings.HasPrefix(h.Liveness, "/") {
		h.Liveness = fmt.Sprintf("/%s", h.Liveness)
	}

	if len(h.Readiness) == 0 {
		h.Readiness = "/ready"
	} else if !strings.HasPrefix(h.Readiness, "/") {
		h.Readiness = fmt.Sprintf("/%s", h.Readiness)
	}

	if len(h.Port) == 0 {
		h.Port = "8080"
	}

	http.Handle(h.Liveness, h.liveness())
	http.Handle(h.Readiness, h.readiness())
	return http.ListenAndServe(fmt.Sprintf(":%s", h.Port), nil)
}

// Ready sets the state of service to ready. State's default value is false.
// You have to manually enabled whenever app is ready to service requests.
func (h *Check) Ready() {
	h.state = ready
}

// NotReady sets the state to notready.
func (h *Check) NotReady() {
	h.state = notready
}

// Terminating starts a go routine waiting for SIGINT & SIGTERM signals.
// Returns true when Kubernetes sends a termination signal to the pod.
// It is a blocking function.
func (h *Check) Terminating() bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()
	return <-done
}

func (h *Check) liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Check) readiness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h.state == ready {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
