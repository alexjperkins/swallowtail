package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/handler"
	accountproto "swallowtail/s.account/proto"
)

const (
	svcName = "s.account"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Init Dao
	if err := dao.Init(ctx, svcName); err != nil {
		panic(err)
	}

	// Mariana Server
	srv := mariana.Init(svcName)
	accountproto.RegisterAccountServer(srv.Grpc(), &handler.AccountService{})
	srv.Run(ctx)
}
