package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.bitfinex/client"
	"swallowtail/s.bitfinex/handler"
	bitfinexproto "swallowtail/s.bitfinex/proto"
)

const (
	svcName = "s.bitfinex"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Init Bitfinex client.
	if err := client.Init(ctx); err != nil {
		panic(err)
	}

	// Init Mariana Server
	srv := mariana.Init(svcName)
	bitfinexproto.RegisterBitfinexServer(srv.Grpc(), &handler.BitfinexService{})
	srv.Run(ctx)
}
