// Package healthz is a tiny & simple to use library for liveness & readiness Kubernetes checks (gRPC included).
package healthz

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
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

// NewCheck .
func NewCheck(live, ready, port string) *Check {
	if len(live) == 0 {
		live = "/live"
	} else if !strings.HasPrefix(live, "/") {
		live = fmt.Sprintf("/%s", live)
	}

	if len(ready) == 0 {
		ready = "/ready"
	} else if !strings.HasPrefix(ready, "/") {
		ready = fmt.Sprintf("/%s", ready)
	}

	if len(port) == 0 {
		port = "8080"
	} else if strings.HasPrefix(port, ":") {
		port = strings.TrimPrefix(port, ":")
	}

	return &Check{
		liveness:  live,
		readiness: ready,
		port:      port,
	}
}

// Start starts the healthcheck http server. It should be called at the start of your application.
// It is a blocking function.
func (h *Check) Start() error {
	srv := &http.Server{
		Handler: h.router(),
		Addr:    fmt.Sprintf(":%s", h.port),
	}
	return srv.ListenAndServe()

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

func (h *Check) router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(h.liveness, h.live).Methods("GET")
	r.HandleFunc(h.readiness, h.ready).Methods("GET")
	return r
}

func (h *Check) live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Check) ready(w http.ResponseWriter, r *http.Request) {
	if h.state == ready {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}
