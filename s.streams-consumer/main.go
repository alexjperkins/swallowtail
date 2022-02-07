package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.streams-consumer/handler"
	streamsconsumerproto "swallowtail/s.streams-consumer/proto"
)

const (
	svcName = "s.streams-consumer"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	srv := mariana.Init(svcName)
	streamsconsumerproto.RegisterStreamsconsumerServer(srv.Grpc(), &handler.StreamsConsumerService{})
	srv.Run(ctx)
}
