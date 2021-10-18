package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.market-data/handler"
	marketdataproto "swallowtail/s.market-data/proto"
)

const (
	svcName = "s.marketdata"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Init Mariana Server
	srv := mariana.Init(svcName)
	marketdataproto.RegisterMarketdataServer(srv.Grpc(), &handler.MarketDataService{})
	srv.Run(ctx)
}
