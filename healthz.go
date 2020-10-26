// Package healthz is a small & simple to use library for liveness & readiness Kubernetes checks (gRPC included).
package healthz

import (
	"context"
	"crypto/x509"
	"log"
	"net/http"
	"time"

	"github.com/jepanetwork/grpc/proto"
	"github.com/johnsiilver/getcert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func Yo() {

}

// LivenessReadiness .
func LivenessReadiness(req chan State, rep chan bool, f func()) {
	http.Handle("/ready", ready(req, rep))
	http.Handle("/live", live(req, rep))
	liserv := func() {
		if err := http.ListenAndServe(":6080", nil); err != nil {
			return
		}
	}

	go liserv()
	go f()
}

func ready(req chan State, rep chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		req <- READY
		switch <-rep {
		case false:
			w.WriteHeader(http.StatusServiceUnavailable)
		case true:
			w.WriteHeader(http.StatusOK)
		}
	}
}

func live(req chan State, rep chan bool) http.HandlerFunc {
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
		switch <-rep {
		case false:
			w.WriteHeader(http.StatusServiceUnavailable)
		case true:
			w.WriteHeader(http.StatusOK)
		}
	}
}

func GRPCLiveCheck() int {
	_, xCert, err := getcert.FromTLSServer("service.jepa.network:50051", true)
	roots := x509.NewCertPool()
	roots.AddCert(xCert[0])
	creds := credentials.NewClientTLSFromCert(roots, "")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial("service.jepa.network:50051", opts...)
	if err != nil {
		log.Println("gRPC connection error: ", err)
	}
	defer conn.Close()

	client := proto.NewHealthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	check, err := client.Check(ctx, &proto.HealthCheckRequest{})
	if err != nil {
		log.Println("Client error: ", err)
	}
	ok := check.GetStatus()
	return int(ok)
}
