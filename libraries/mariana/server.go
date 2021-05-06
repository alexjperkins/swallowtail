package mariana

import (
	"net"

	"google.golang.org/grpc"
)

type Server interface {
	Serve(listener net.Listener) error
}

func NewServer() Server {
	var opts []grpc.ServerOption
	gs := grpc.NewServer(opts...)
	return &server{
		gs: gs,
	}
}

type server struct {
	// GRPC server.
	gs *grpc.Server
}

func (g *server) Serve(listener net.Listener) error {
	return g.gs.Serve(listener)
}
