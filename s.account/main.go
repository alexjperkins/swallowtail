package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"swallowtail/libraries/util"
	"swallowtail/s.account/dao"
	"syscall"

	"github.com/jackc/pgconn"
)

var (
	pgUser     string
	pgPassword string
	pgHost     string
	pgPort     string
	pgDB       string
)

func init() {
	pgUser = util.SetEnv("POSTGRES_USER")
	pgPassword = util.SetEnv("POSTGRES_PASSWORD")
	pgHost = util.SetEnv("POSTGRES_HOST")
	pgPort = util.SetEnv("POSTGRES_PORT")
	pgDB = util.SetEnv("POSTGRES_DB")
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	port, err := strconv.ParseUint(pgPort, 10, 16)
	if err != nil {
		panic(fmt.Sprintf("Cannot parse port: %v; %v", pgPort, err))
	}

	// Options
	opts := &pgconn.Config{
		Host:     pgHost,
		Port:     uint16(port),
		Database: pgDB,
		User:     pgUser,
		Password: pgPassword,
	}

	dbCloser, err := dao.Init(ctx, opts)
	if err != nil {
		panic(err)
	}

	defer dbCloser()
	defer cancel()

	// Only whilst we don't have a server running.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	select {
	case <-sc:
		return
	}
}
