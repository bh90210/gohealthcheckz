// Package healthz is a small & simple to use library for liveness & readiness Kubernetes checks (gRPC included).
package healthz

import (
	"net/http"
)

type Service int

const (
	LOGGER Service = iota
	GRPC
	AUTHORIZATION
	AUTHENTICATION
	WEBSITE
)

type State int

const (
	READY State = iota
	LIVE
)

// LivenessReadiness .
func LivenessReadiness(req chan State, res chan bool, f func()) {
	http.Handle("/ready", ready(req, res))
	http.Handle("/live", live(req, res))
	liserv := func() {
		if err := http.ListenAndServe(":6080", nil); err != nil {
			return
		}
	}

	go liserv()
	go f()
}

func ready(req chan State, res chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		req <- READY
		switch <-res {
		case false:
			w.WriteHeader(http.StatusServiceUnavailable)
		case true:
			w.WriteHeader(http.StatusOK)
		}
	}
}

func live(req chan State, res chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		// switch *service {
		// case GRPC:
		// 	switch gRPCLiveCheck() {
		// 	case 1:
		// 		w.WriteHeader(http.StatusOK)
		// 	default:
		// 		w.WriteHeader(http.StatusServiceUnavailable)
		// 	}
		// case LOGGER:
		// 	w.WriteHeader(http.StatusOK)
		// }
		req <- LIVE
		switch <-res {
		case false:
			w.WriteHeader(http.StatusServiceUnavailable)
		case true:
			w.WriteHeader(http.StatusOK)
		}
	}
}

// startup

// termination
