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

type state bool

const (
	notready state = false
	ready    state = true
)

type Check struct {
	liveness  string
	readiness string
	port      string
	state     state
}

func NewCheck(live, ready, port string) *Check {
	h := &Check{}

	if len(h.liveness) == 0 {
		h.liveness = "/live"
	} else if !strings.HasPrefix(h.liveness, "/") {
		h.liveness = fmt.Sprintf("/%s", h.liveness)
	}

	if len(h.readiness) == 0 {
		h.readiness = "/ready"
	} else if !strings.HasPrefix(h.readiness, "/") {
		h.readiness = fmt.Sprintf("/%s", h.readiness)
	}

	if len(h.port) == 0 {
		h.port = "8080"
	} else if strings.HasPrefix(h.readiness, ":") {
		h.readiness = strings.TrimPrefix(h.readiness, ":")
	}

	return h
}

// Start starts the healthcheck http server. It should be called at the start of your application.
// It is a blocking function.
func (h *Check) Start() error {
	http.Handle(h.liveness, h.live())
	http.Handle(h.readiness, h.ready())
	return http.ListenAndServe(fmt.Sprintf(":%s", h.port), nil)
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

func (h *Check) live() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Check) ready() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h.state == ready {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
