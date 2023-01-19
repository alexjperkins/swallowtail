package main

import (
	"context"
	"log"

	"swallowtail/libraries/environment"
	"swallowtail/libraries/mariana"
	"swallowtail/s.bookmarker/handler"
	bookmarkerproto "swallowtail/s.bookmarker/proto"
)

const serviceName = "s.bookmarker"

func main() {
	cfg, err := environment.LoadEnvironment()
	if err != nil {
		log.Fatalf("Failed to load environment: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init Mariana Server.
	srv := mariana.InitWithConfig(serviceName, cfg)
	bookmarkerproto.RegisterBookmarkerServer(srv.Grpc(), &handler.BookmarkerService{})

	// Run server.
	srv.Run(ctx)
}
