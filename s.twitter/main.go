package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"swallowtail/s.twitter/consumers"
	"syscall"

	"github.com/monzo/slog"
)

func main() {
	ctx := context.Background()
	tw := consumers.New()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	defer slog.Warn(ctx, "Received shutdown signal....")

	slog.Info(ctx, "Starting twitter consumer...")
	stopFunc, err := tw.Run(ctx)
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}

	select {
	case <-sc:
		stopFunc()
		tw.Done(ctx)

	case <-ctx.Done():
		stopFunc()
		tw.Done(ctx)
	}
}
