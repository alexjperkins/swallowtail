package main

import (
	"context"
	"swallowtail/libraries/mariana"
	"swallowtail/s.payments/dao"
	"swallowtail/s.payments/handler"
	paymentsproto "swallowtail/s.payments/proto"
)

const (
	svcName = "s.payments"
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
	paymentsproto.RegisterPaymentsServer(srv.Grpc(), &handler.PaymentsService{})
	srv.Run(ctx)
}
