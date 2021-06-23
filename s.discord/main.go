package main

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/jackc/pgconn"
	"google.golang.org/grpc"

	"swallowtail/libraries/util"
	"swallowtail/s.discord/dao"
	"swallowtail/s.discord/handler"
	discordproto "swallowtail/s.discord/proto"
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

	// Dao
	port, err := strconv.ParseUint(pgPort, 10, 16)
	if err != nil {
		panic(fmt.Sprintf("Cannot parse port: %v; %v", pgPort, err))
	}

	daoOpts := &pgconn.Config{
		Host:     pgHost,
		Port:     uint16(port),
		Database: pgDB,
		User:     pgUser,
		Password: pgPassword,
	}

	dbCloser, err := dao.Init(ctx, daoOpts)
	if err != nil {
		panic(err)
	}

	defer dbCloser()
	defer cancel()

	// gRPC server
	lis, err := net.Listen("tcp", "8000")
	if err != nil {
		panic(nil)
	}

	s := grpc.NewServer()
	discordproto.RegisterDiscordServer(s, &handler.DiscordService{})
	if err := s.Serve(lis); err != nil {
	}
}
