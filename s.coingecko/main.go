package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.coingecko/client"
	"swallowtail/s.coingecko/handler"
	coingeckoproto "swallowtail/s.coingecko/proto"
)

const (
	svcName = "s.coingecko"
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
	coingeckoproto.RegisterCoingeckoServer(srv.Grpc(), &handler.CoingeckoService{})

	srv.Run(ctx)
}
