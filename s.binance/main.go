package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.binance/client"
	"swallowtail/s.binance/exchangeinfo"
	"swallowtail/s.binance/handler"
	binanceproto "swallowtail/s.binance/proto"

	"github.com/monzo/slog"
)

const (
	svcName = "s.binance"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Init Binance client.
	if err := client.Init(ctx); err != nil {
		panic(err)
	}

	slog.Info(ctx, "HERE")

	// Init exchange info.
	if err := exchangeinfo.Init(ctx); err != nil {
		panic(err)
	}

	// Init Mariana Server
	srv := mariana.Init(svcName)
	binanceproto.RegisterBinanceServer(srv.Grpc(), &handler.BinanceService{})
	srv.Run(ctx)
}
