package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.trade-engine/dao"
	"swallowtail/s.trade-engine/handler"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

const (
	svcName = "s.trade-engine"
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
	tradeengineproto.RegisterTradeengineServer(srv.Grpc(), &handler.TradeEngineService{})
	srv.Run(ctx)
}
