package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.bookmarker/handler"
	bookmarkerproto "swallowtail/s.bookmarker/proto"
)

const serviceName = "s.bookmarker"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init Mariana Server.
	srv := mariana.Init(serviceName)
	bookmarkerproto.RegisterBookmarkerServer(srv.Grpc(), &handler.BookmarkerService{})

	// Run server.
	srv.Run(ctx)
}
