package grpc

import (
	"context"
	"crypto/x509"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/johnsiilver/getcert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

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
