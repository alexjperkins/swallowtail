package main

import (
	"context"
	"os"
	"os/signal"
	"swallowtail/s.binance/client"
	"syscall"

	"github.com/monzo/slog"
)

func main() {
	ctx := context.Background()
	client.Init()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	defer slog.Warn(ctx, "Received shutdown signal....")

	select {
	case <-sc:
	}
}
