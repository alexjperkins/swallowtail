package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.market-data/handler"
	marketdataproto "swallowtail/s.market-data/proto"
	"swallowtail/s.web3/client"
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
	marketdataproto.RegisterMarketdataServer(srv.Grpc(), &handler.MarketDataService{})
	srv.Run(ctx)
}
