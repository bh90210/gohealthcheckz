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

type Check struct {
	liveness  string
	readiness string
	port      string
	ready     bool
}

// NewCheck initializes and return a new Check.
// Default values for server: endpoints (/live, /ready) and port `:8080`.
func NewCheck(options ...func(*Check)) *Check {
	// set defaults
	c := &Check{liveness: "/live", readiness: "/ready", port: ":8080"}
	// range over provided options to overwrite defaults
	for _, option := range options {
		option(c)
	}
	return c
}

// Start starts the healthcheck http server. It should be called at the start of your application.
// It is a blocking function.
func (h *Check) Start() error {
	srv := &http.Server{
		Handler: h.router(),
		Addr:    h.port,
	}
	return srv.ListenAndServe()
}

// Ready sets the state of service to ready. State's default value is false.
// You have to manually enabled whenever app is ready to service requests.
func (h *Check) Ready() {
	h.ready = true
}

// NotReady sets the state to notready.
func (h *Check) NotReady() {
	h.ready = false
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

// OptionsLivePath sets live path endpoint.
func OptionsLivePath(path string) func(*Check) {
	return func(c *Check) {
		if !strings.HasPrefix(path, "/") {
			path = fmt.Sprintf("/%s", path)
		}
		c.liveness = path
	}
}

// OptionsReadyPath sets ready path endpoint.
func OptionsReadyPath(path string) func(*Check) {
	return func(c *Check) {
		if !strings.HasPrefix(path, "/") {
			path = fmt.Sprintf("/%s", path)
		}
		c.readiness = path
	}
}

// OptionsPort sets health check server's port.
func OptionsPort(port string) func(*Check) {
	return func(c *Check) {
		if !strings.HasPrefix(port, ":") {
			port = fmt.Sprintf(":%s", port)
		}
		c.port = port
	}
}

func (h *Check) router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(h.liveness, h.liveHandler).Methods("GET")
	r.HandleFunc(h.readiness, h.readyHandler).Methods("GET")
	return r
}

func (h *Check) liveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Check) readyHandler(w http.ResponseWriter, r *http.Request) {
	if h.ready {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}
