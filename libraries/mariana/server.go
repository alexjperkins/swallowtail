package mariana

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/monzo/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"swallowtail/libraries/environment"
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
// TODO: deprecate.
func Init(serviceName string) Server {
	return initServer(serviceName, nil)
}

// InitWithConfig ...
func InitWithConfig(serviceName string, cfg *environment.Environment) Server {
	return initServer(serviceName, cfg)
}

func initServer(serviceName string, cfg *environment.Environment) *server {
	grpcs := grpc.NewServer()

	reflection.Register(grpcs)
	s := &server{
		s:           grpcs,
		ServiceName: serviceName,
	}

	if cfg != nil {
		s.Config = cfg
	}

	return s
}

type server struct {
	// Service Name
	ServiceName string
	Config      *environment.Environment

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
