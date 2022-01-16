package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.ftx/client"
	"swallowtail/s.ftx/exchangeinfo"
	"swallowtail/s.ftx/handler"
	ftxproto "swallowtail/s.ftx/proto"
)

const (
	svcName = "s.ftx"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Init client.
	if err := client.Init(ctx); err != nil {
		panic(err)
	}

	// Init exchangeinfo.
	if err := exchangeinfo.Init(ctx); err != nil {
		panic(err)
	}

	// Init Mariana Server
	srv := mariana.Init(svcName)
	ftxproto.RegisterFtxServer(srv.Grpc(), &handler.FTXService{})
	srv.Run(ctx)
}
