package client

import (
	"context"
	"strconv"
	"swallowtail/s.googlesheets/templates"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"google.golang.org/api/sheets/v4"
)

type googlesheetsClient struct {
	s *sheets.Service
}

func (g *googlesheetsClient) Ping(ctx context.Context) bool {
	return false
}

func (g *googlesheetsClient) Values(_ context.Context, sheetID string, rowsRange string) ([][]interface{}, error) {
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

func (g *googlesheetsClient) CreateSheet(ctx context.Context, sheetType templates.SheetType) (string, error) {
	// Create the initial spreadsheet.
	call := g.s.Spreadsheets.Create(&sheets.Spreadsheet{})
	call = call.Context(ctx)
	s, err := call.Do()
	if err != nil {
		return "", terrors.Augment(err, "Failed to create a googlesheet", nil)
	}

	template, err := templates.GetTemplateByType(sheetType)
	switch {
	case terrors.Is(err, "template-does-not-exist"):
		// Nothing more to do here
		return s.SpreadsheetUrl, nil
	}

	if err := g.UpdateRows(ctx, s.SpreadsheetId, template.RowRange(), template.Values()); err != nil {
		return "", terrors.Augment(err, "Failed to update rows for template", map[string]string{
			"template": template.ID().String(),
		})
	}

	return s.SpreadsheetUrl, nil
}
