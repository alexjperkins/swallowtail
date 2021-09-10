package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.satoshi/handler"
	"swallowtail/s.satoshi/parser"
	satoshiproto "swallowtail/s.satoshi/proto"
	"swallowtail/s.satoshi/satoshi"
)

const (
	svcName = "s.satoshi"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Init parser.
	if err := parser.Init(ctx); err != nil {
		panic(err)
	}

	// Init our background satoshi jobs.
	if err := satoshi.Init(ctx); err != nil {
		panic(err)
	}

	// Init Mariana Server.
	srv := mariana.Init(svcName)
	satoshiproto.RegisterSatoshiServer(srv.Grpc(), &handler.SatoshiService{})
	srv.Run(ctx)
}
