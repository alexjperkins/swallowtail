package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"swallowtail/libraries/util"
	"swallowtail/s.discord/handler"
	"syscall"
	"time"

	"github.com/monzo/slog"
	"github.com/monzo/typhon"
)

func main() {
	hostname := util.SetEnv("S_DISCORD_HOSTNAME")
	port := util.SetEnv("S_DISCORD_PORT")

	addr, err := util.ServiceRPCAddress(hostname, port)
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))

	}

	svc := handler.Service().
		Filter(typhon.ErrorFilter).
		Filter(typhon.H2cFilter)
	srv, err := typhon.Listen(svc, addr)
	if err != nil {
		panic(err)
	}

	slog.Info(nil, "👋  Listening on %v", srv.Listener().Addr())

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	slog.Warn(nil, "☠️  Shutting down")
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Stop(c)
}
