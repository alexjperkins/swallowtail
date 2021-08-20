package client

import (
	"context"
	"fmt"
	"strconv"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.googlesheets/templates"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

type googlesheetsClient struct {
	s *sheets.Service
	// We need this to handle permissions; super annoying.
	d *drive.Service
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
		return terrors.Augment(err, "Failed to update rows", map[string]string{
			"row_range": rowsRange,
		})
	}
	return nil
}

func (g *googlesheetsClient) CreateSheet(ctx context.Context, sheetType templates.SheetType, emailAddress string) (*sheets.Spreadsheet, error) {
	errParams := map[string]string{
		"email_address": emailAddress,
	}

	// Create the initial spreadsheet.
	create := g.s.Spreadsheets.Create(&sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: fmt.Sprintf("Swallowtail Portfolio: %s", emailAddress),
		},
		Sheets: []*sheets.Sheet{
			{
				Properties: &sheets.SheetProperties{
					Title: "SpotPortfolio",
				},
			},
		},
	})
	create = create.Context(ctx)
	s, err := create.Do()
	if err != nil {
		return nil, terrors.Augment(err, "Failed to create a googlesheet", errParams)
	}

	slog.Info(ctx, "Created sheet: %s, with metadata: %v", s.SpreadsheetUrl, s.DeveloperMetadata)

	// Give permissions via the drive api.
	givePermissions := g.d.Permissions.Create(s.SpreadsheetId, &drive.Permission{
		EmailAddress: emailAddress,
		Role:         "writer",
		Type:         "user",
	})
	givePermissions.Context(ctx)

	if _, err := givePermissions.Do(); err != nil {
		return nil, terrors.Augment(err, "Failed to create a googlesheet", errParams)
	}

	// Write template to the file.
	template, err := templates.GetTemplateByType(sheetType)
	switch {
	case terrors.Is(err, "template-does-not-exist"):
		// Nothing more to do here
		return s, nil
	}

	if err := g.UpdateRows(ctx, s.SpreadsheetId, template.RowRange(), template.Values()); err != nil {
		return nil, terrors.Augment(err, "Failed to update rows for template", map[string]string{
			"email_address": emailAddress,
			"template":      template.ID().String(),
		})
	}

	return s, nil
}

func (g *googlesheetsClient) RegisterSheet(ctx context.Context, spreadsheetID, emailAddress string) error {
	errParams := map[string]string{
		"spreadsheet_id": spreadsheetID,
		"email_address":  emailAddress,
	}

	// Give permissions via the drive api.
	givePermissions := g.d.Permissions.Create(spreadsheetID, &drive.Permission{
		EmailAddress: emailAddress,
		Role:         "writer",
		Type:         "user",
	})
	givePermissions.Context(ctx)

	if _, err := givePermissions.Do(); err != nil {
		return gerrors.Augment(err, "Failed to create a googlesheet", errParams)
	}

	slog.Info(ctx, "Gave %s %s permissions", emailAddress, spreadsheetID)

	return nil
}
