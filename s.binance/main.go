package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.account/dao"
	"swallowtail/s.binance/client"
	"swallowtail/s.binance/handler"
	binanceproto "swallowtail/s.binance/proto"
)

const (
	svcName = "s.binance"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Init Dao
	if err := dao.Init(ctx, svcName); err != nil {
		panic(err)
	}

	// Init Mariana Server
	srv := mariana.Init(svcName)
	binanceproto.RegisterBinanceServer(srv.Grpc(), &handler.BinanceService{})
	srv.Run(ctx)

	// Init Binance client.
	if err := client.Init(ctx); err != nil {
		panic(err)
	}
}
