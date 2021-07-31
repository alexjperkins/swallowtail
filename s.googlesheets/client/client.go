package client

import (
	"context"
	"io/ioutil"

	"github.com/monzo/terrors"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	"swallowtail/s.googlesheets/templates"
)

var (
	// TODO refactor
	client               GooglesheetsClient
	defaultConfigPath    = "/home/alexjperkins/repos/swallowtail/s.googlesheets/clients/credentials.json"
	defaultTokenFilePath = "./token.json"
)

// GooglesheetsClient ...
type GooglesheetsClient interface {
	Ping(ctx context.Context) bool
	Values(ctx context.Context, sheetID string, rowsRange string) ([][]interface{}, error)
	UpdateRows(ctx context.Context, sheetID, rowsRange string, values [][]interface{}) error
	CreateSheet(ctx context.Context, sheetType templates.SheetType) (string, error)
}

// Init sets a new google sheets client for interacting with googlesheets.
func Init(ctx context.Context) error {
	// TODO; we need to refactor this.
	b, err := ioutil.ReadFile(defaultConfigPath)
	if err != nil {
		return terrors.Augment(err, "Failed to load googlesheets credentials", nil)
	}
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return terrors.Augment(err, "Failed to load config", nil)
	}
	c := getClient(config)
	srv, err := sheets.New(c)
	if err != nil {
		return terrors.Augment(err, "Failed to create sheets client", nil)
	}

	client = &googlesheetsClient{
		s: srv,
	}

	return nil
}

// Values ...
func Values(ctx context.Context, sheetID string, rowsRange string) ([][]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Retrieve values from googlesheets spreadsheet")
	defer span.Finish()
	return client.Values(ctx, sheetID, rowsRange)
}

// UpdateRows ...
func UpdateRows(ctx context.Context, sheetID string, rowsRange string, values [][]interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Update values from googlesheets spreadsheet")
	defer span.Finish()
	return client.UpdateRows(ctx, sheetID, rowsRange, values)
}

// CreateSheet ...
func CreateSheet(ctx context.Context, sheetType templates.SheetType) (string, error) {
	// TODO: metrics & opnetracing.
	return client.CreateSheet(ctx, sheetType)
}

func isValidHTTPStatusCode(c int) bool {
	if c < 200 || c > 299 {
		return false
	}
	return true
}
