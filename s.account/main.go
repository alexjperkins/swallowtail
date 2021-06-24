package main

import (
	"context"
	"net"

	"swallowtail/s.account/dao"
	"swallowtail/s.account/handler"
	accountproto "swallowtail/s.account/proto"

	"google.golang.org/grpc"
)

const (
	serviceName = "s.account"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	// Dao
	err := dao.Init(ctx, serviceName)
	if err != nil {
		panic(err)
	}

	// gRPC server
	lis, err := net.Listen("tcp", "8000")
	if err != nil {
		panic(nil)
	}

	s := grpc.NewServer()
	accountproto.RegisterAccountServer(s, &handler.AccountService{})
	if err := s.Serve(lis); err != nil {
		cancel()
	}
}
