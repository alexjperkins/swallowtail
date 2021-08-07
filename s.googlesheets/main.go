package main

import (
	"context"
	"swallowtail/libraries/mariana"
	"swallowtail/s.googlesheets/client"
	"swallowtail/s.googlesheets/dao"
	"swallowtail/s.googlesheets/handler"
	googlesheetsproto "swallowtail/s.googlesheets/proto"
	"swallowtail/s.googlesheets/sync"
)

var (
	// Move
	// defaultAlexGoogleSpreadsheetID = "1AYtRsdEcoEjmh-OtribxJ9et7qvCf6Z_UkkYNnKqqZY"
	// defaultBenGoogleSpreadsheetID  = "1Krg7O8h-ItK42dTC-ey9HOh6v1w8T2SCHcUsVJdUnmI"

	svcName = "s.googlesheets"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := dao.Init(ctx, svcName); err != nil {
		panic(err)
	}

	// Init googlesheets client.
	if err := client.Init(ctx); err != nil {
		panic(err)
	}

	// Init syncer.
	if err := sync.Init(ctx); err != nil {
		panic(err)
	}

	// Init Mariana Server
	srv := mariana.Init(svcName)
	googlesheetsproto.RegisterGooglesheetsServer(srv.Grpc(), &handler.GooglesheetsService{})
	srv.Run(ctx)
}
