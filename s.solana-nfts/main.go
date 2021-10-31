package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.solana-nfts/client"
	"swallowtail/s.solana-nfts/handler"
	solananftsproto "swallowtail/s.solana-nfts/proto"
)

const (
	svcName = "s.solananfts"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Init vendor clients.
	if err := client.Init(ctx); err != nil {
		panic(err)
	}

	// Init Mariana Server
	srv := mariana.Init(svcName)
	solananftsproto.RegisterSolananftsServer(srv.Grpc(), &handler.SolanaNFTsService{})
	srv.Run(ctx)
}
