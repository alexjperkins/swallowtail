package client

import (
	"context"
	"io/ioutil"

	"github.com/monzo/terrors"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"

	"swallowtail/s.googlesheets/templates"
)

const (
	// ScopeSpreadsheets the scope we use for our spreadsheets
	ScopeSpreadsheets = "https://www.googleapis.com/auth/spreadsheets"
	ScopeDrive        = "https://www.googleapis.com/auth/drive"
)

var (
	client GooglesheetsClient
)

// GooglesheetsClient ...
type GooglesheetsClient interface {
	Ping(ctx context.Context) bool
	Values(ctx context.Context, sheetID string, rowsRange string) ([][]interface{}, error)
	UpdateRows(ctx context.Context, sheetID, rowsRange string, values [][]interface{}) error
	CreateSheet(ctx context.Context, sheetType templates.SheetType, emailAddress string) (*sheets.Spreadsheet, error)
}

// Init sets a new google sheets client for interacting with googlesheets.
func Init(ctx context.Context) error {
	jwtJSON, err := ioutil.ReadFile("/s.googlesheets/config/credentials.json")
	if err != nil {
		return terrors.Augment(err, "Failed to read credentials", nil)
	}

	cfg, err := google.JWTConfigFromJSON(jwtJSON, ScopeSpreadsheets, ScopeDrive)
	if err != nil {
		return terrors.Augment(err, "Failed to create JWT config", nil)
	}

	c := cfg.Client(ctx)

	s, err := sheets.New(c)
	if err != nil {
		return terrors.Augment(err, "Failed to create google sheets client", nil)
	}

	d, err := drive.New(c)
	if err != nil {
		return terrors.Augment(err, "Failed to create google drive client", nil)
	}

	client = &googlesheetsClient{
		s: s,
		d: d,
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
func CreateSheet(ctx context.Context, sheetType templates.SheetType, emailAddress string) (*sheets.Spreadsheet, error) {
	// TODO: metrics & opnetracing.
	return client.CreateSheet(ctx, sheetType, emailAddress)
}

func isValidHTTPStatusCode(c int) bool {
	if c < 200 || c > 299 {
		return false
	}
	return true
}
