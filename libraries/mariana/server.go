package mariana

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/monzo/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	network            = "tcp"
	defaultServicePort = "8000"
)

// Server defines the interface for our base server setup.
type Server interface {
	Run(ctx context.Context)
	Grpc() *grpc.Server
}

// Init inits our base server.
func Init(service string) Server {
	s := grpc.NewServer()
	reflection.Register(s)
	return &server{
		s:           s,
		ServiceName: service,
	}
}

type server struct {
	// Service Name
	ServiceName string

	// GRPC server.
	s *grpc.Server
}

// Run runs our base server.
func (s *server) Run(ctx context.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("Failed to establish hostname: %v", err))
	}

	addr := formatAddr(hostname, defaultServicePort)
	errParams := map[string]string{
		"addr":         addr,
		"service_name": s.ServiceName,
	}

	listener, err := net.Listen(network, addr)
	if err != nil {
		panic(fmt.Sprintf("%s failed to listen on %s:%s: %v", s.ServiceName, network, addr, err))
	}

	slog.Info(ctx, "%s listening on %s", s.ServiceName, addr, errParams)

	if err := s.s.Serve(listener); err != nil {
		panic(err)
	}
}

// Grpc returns the underlying gRPC server.
func (s *server) Grpc() *grpc.Server {
	return s.s
}

func formatAddr(serviceName, port string) string {
	return fmt.Sprintf("%s:%s", serviceName, port)
}
