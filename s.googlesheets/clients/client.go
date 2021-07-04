package clients

import (
	"context"
	"io/ioutil"
	"strconv"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var (
	defaultConfigPath    = "/home/alexjperkins/repos/swallowtail/s.googlesheets/clients/credentials.json"
	defaultTokenFilePath = "./token.json"
)

type GooglesheetsClient interface {
	Ping(ctx context.Context) bool
	Values(sheetID string, rowsRange string) ([][]interface{}, error)
	UpdateRows(ctx context.Context, sheetID, rowsRange string, values [][]interface{}) error
}

type googlesheetsClient struct {
	s *sheets.Service
}

// New creates a new google sheets client for interacting with googlesheets.
func New(ctx context.Context) (*googlesheetsClient, error) {
	b, err := ioutil.ReadFile(defaultConfigPath)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to load googlesheets credentials", nil)
	}
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, terrors.Augment(err, "Failed to load config", nil)
	}
	client := getClient(config)
	srv, err := sheets.New(client)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to create sheets client", nil)
	}

	return &googlesheetsClient{
		s: srv,
	}, nil
}

func (g *googlesheetsClient) Ping(ctx context.Context) bool {
	return false
}

func (g *googlesheetsClient) Values(sheetID string, rowsRange string) ([][]interface{}, error) {
	r, err := g.s.Spreadsheets.Values.Get(sheetID, rowsRange).Do()
	if err != nil {
		return nil, err
	}
	if !isValidHTTPStatusCode(r.HTTPStatusCode) {
		return nil, terrors.BadRequest("spreadsheet-load-failed", "Failed to load spreadsheet", map[string]string{
			"spreadsheetID":    sheetID,
			"http_status_code": strconv.Itoa(r.HTTPStatusCode),
		})
	}
	return r.Values, nil
}

func (g *googlesheetsClient) UpdateRows(ctx context.Context, sheetID, rowsRange string, values [][]interface{}) error {
	context.Background()

	v := &sheets.ValueRange{
		Range:  rowsRange,
		Values: values,
	}
	req := g.s.Spreadsheets.Values.Update(sheetID, rowsRange, v)
	req.ValueInputOption("RAW")

	if _, err := req.Do(); err != nil {
		slog.Error(ctx, err.Error())
		return terrors.Augment(err, "Failed to update rows", map[string]string{
			"row_range": rowsRange,
		})
	}
	return nil
}

func isValidHTTPStatusCode(c int) bool {
	if c < 200 || c > 299 {
		return false
	}
	return true
}
