package mariana

import (
	"context"
	"fmt"
	"net"

	"github.com/monzo/slog"
	"google.golang.org/grpc"
)

const (
	network            = "tcp"
	defaultServicePort = "8000"
)

type Server interface {
	Run(ctx context.Context)
	Grpc() *grpc.Server
}

func Init(service string) Server {
	s := grpc.NewServer()
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

func (s *server) Run(ctx context.Context) {
	addr := formatAddr(s.ServiceName, defaultServicePort)
	errParams := map[string]string{
		"addr":         addr,
		"service_name": s.ServiceName,
	}

	listener, err := net.Listen(network, addr)
	if err != nil {
		panic(err)
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
