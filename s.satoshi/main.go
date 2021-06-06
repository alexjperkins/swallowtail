package main

import (
	"context"
	"os"
	"os/signal"
	"swallowtail/s.satoshi/satoshi"
	"syscall"

	"github.com/monzo/slog"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	satoshi := satoshi.New(true)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	defer slog.Warn(ctx, "Received shutdown signal....")

	slog.Info(ctx, "Starting Satoshi...")
	satoshi.Run(ctx)
	select {
	case <-sc:
		satoshi.Stop()
		cancel()
		return
	}
}
