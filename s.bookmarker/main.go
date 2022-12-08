package main

import (
	"context"
	"log"

	"swallowtail/libraries/environment"
	"swallowtail/libraries/mariana"
	"swallowtail/s.bookmarker/dao"
	"swallowtail/s.bookmarker/handler"
	bookmarkerproto "swallowtail/s.bookmarker/proto"

	"github.com/monzo/slog"
)

func main() {
	// Load environment.
	cfg, err := environment.LoadEnvironment()
	if err != nil {
		log.Fatalf("Failed to load environment: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init Mariana Server.
	srv := mariana.Init(cfg.Metadata.ServiceName)
	bookmarkerproto.RegisterBookmarkerServer(srv.Grpc(), &handler.BookmarkerService{})

	// Init Dao.
	if err := dao.Init(cfg.Cassandra); err != nil {
        // TODO: lets panic here.
		slog.Error(ctx, "Failed to initialize dao: %v", err)
	}

	// Run server.
	srv.Run(ctx)
}
