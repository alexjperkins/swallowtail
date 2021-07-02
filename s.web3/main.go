package main

import (
	"context"
	"os"
	"os/signal"
	"swallowtail/s.web3/client"
	"syscall"

	"github.com/monzo/slog"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	client.Init(ctx)

	select {
	case <-sc:
		slog.Warn(ctx, "Received shutdown signal....")
		cancel()
		return
	}
}
