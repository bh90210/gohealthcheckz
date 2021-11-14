package grpc

import (
	"context"
	"crypto/x509"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/johnsiilver/getcert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	ready bool
}

func NewCheckGRPC(options ...func(*Server)) (*Server, error) {
	s := &Server{}
	for _, option := range options {
		option(s)
	}
	return s, nil
}

func (h *Server) Start() error {
	// srv := &http.Server{
	// 	Handler: h.router(),
	// 	Addr:    fmt.Sprintf(":%s", h.port),
	// }
	// return srv.ListenAndServe()
	return nil
}

// Ready sets the state of service to ready. State's default value is false.
// You have to manually enabled whenever app is ready to service requests.
func (h *Server) Ready() {
	h.ready = true
}

// NotReady sets the state to notready.
func (h *Server) NotReady() {
	h.ready = false
}

// Terminating starts a go routine waiting for SIGINT & SIGTERM signals.
// Returns true when Kubernetes sends a termination signal to the pod.
// It is a blocking function.
func (h *Server) Terminating() bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()
	return <-done
}

func LivePath(path string) func(*Server) {
	return func(s *Server) {
		s.ready = true
	}
}

func ReadyPath(s *Server) *Server {
	return s
}

func Port(s *Server) *Server {
	return s
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
