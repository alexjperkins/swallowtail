package main

import (
	"context"

	"swallowtail/libraries/mariana"
	"swallowtail/s.streams-producer/handler"
	streamsproducerproto "swallowtail/s.streams-producer/proto"
)

const (
	svcName = "s.streams-producer"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	srv := mariana.Init(svcName)
	streamsproducerproto.RegisterStreamsproducerServer(srv.Grpc(), &handler.StreamsProducerService{})
	srv.Run(ctx)
}
