package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.bybt/client"
	"swallowtail/s.bybt/handler"
	bybtproto "swallowtail/s.bybt/proto"
)

const (
	svcName = "s.bybt"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := client.Init(ctx); err != nil {
		panic(err)
	}

	// Init Mariana Server
	srv := mariana.Init(svcName)
	bybtproto.RegisterBybtServer(srv.Grpc(), &handler.ByBtService{})
	srv.Run(ctx)
}
