package spreadsheet

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"swallowtail/s.googlesheets/clients"
	"swallowtail/s.googlesheets/domain"
	"swallowtail/s.googlesheets/owner"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

const (
	defaultPortfolioRowRange          = "C8:L100"
	defaultPortfolioMetatdataRange    = "D2:D4"
	defaultPortfolioTradeHistoryRange = "Q8:W100"
)

func New(id string, sheetIDs []string, client clients.GooglesheetsClient, owner *owner.GooglesheetOwner) *GoogleSheetPortfolio {
	return &GoogleSheetPortfolio{
		Owner:    owner,
		SheetIDs: sheetIDs,
		id:       id,
		c:        client,
	}
}

type GoogleSheetPortfolio struct {
	Owner    *owner.GooglesheetOwner
	SheetIDs []string
	c        clients.GooglesheetsClient
	id       string
}

func (gsp *GoogleSheetPortfolio) ID() string {
	return gsp.id
}

func (gsp *GoogleSheetPortfolio) Rows(ctx context.Context, sheetID string) ([]*domain.PortfolioRow, error) {
	fs := formatRowsRange(sheetID)
	vs, err := gsp.c.Values(gsp.ID(), fs)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to retreive sheet values", map[string]string{
			"sheet_id":   gsp.ID(),
			"owner_name": gsp.Owner.Name,
		})
	}
	rows := []*domain.PortfolioRow{}
	if vs == nil {
		return rows, nil
	}
	for i, v := range vs {
		if len(v) == 0 {
			continue
		}

		r, err := gsp.parseRow(ctx, i, v)
		if err != nil {
			return nil, terrors.Augment(err, "Failed to parse row", map[string]string{
				"sheet_id":   gsp.ID(),
				"owner_name": gsp.Owner.Name,
			})
		}
		rows = append(rows, r)
	}
	slog.Info(ctx, "All rows parsed", map[string]string{
		"sheet_id":    sheetID,
		"owner_name":  gsp.Owner.Name,
		"number_rows": strconv.Itoa(len(rows)),
	})
	return gsp.validateRows(rows)
}

func (gsp *GoogleSheetPortfolio) Metadata(ctx context.Context) (*domain.PortfolioMetadata, error) {
	sheetID := gsp.SheetIDs[0]
	r := formatMetadataRange(sheetID)
	slog.Info(ctx, "Fetching metadata.", map[string]string{
		"owner_name": gsp.Owner.Name,
		"row_range":  r,
	})
	vs, err := gsp.c.Values(gsp.ID(), r)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to retrieve sheet pnl metadata", map[string]string{
			"sheet_id":   gsp.ID(),
			"owner_name": gsp.Owner.Name,
			"range":      r,
		})
	}
	if vs == nil {
		return &domain.PortfolioMetadata{}, nil
	}
	return gsp.parseMetadata(ctx, vs)
}

func (gsp *GoogleSheetPortfolio) TradeHistory(ctx context.Context) ([]*domain.HistoricalTradeRow, error) {
	sheetID := gsp.SheetIDs[0]
	r := formatTradeHistoryRange(sheetID)
	vs, err := gsp.c.Values(gsp.ID(), r)

	slog.Info(ctx, "Fetching trade history.", map[string]string{
		"owner_name": gsp.Owner.Name,
		"row_range":  r,
	})
	if err != nil {
		return nil, terrors.Augment(err, "Failed to retrieve sheet trade history", map[string]string{
			"sheet_id":   gsp.ID(),
			"owner_name": gsp.Owner.Name,
			"range":      r,
		})
	}
	rows := []*domain.HistoricalTradeRow{}
	if vs == nil {
		return rows, nil
	}
	for i, v := range vs {
		if len(v) == 0 {
			continue
		}
		r, err := gsp.parseTradeHistoryRow(ctx, i, v)
		if err != nil {
			return nil, terrors.Augment(err, "Failed to parse historical trade row", map[string]string{
				"sheet_id":   gsp.ID(),
				"owner_name": gsp.Owner.Name,
			})
		}
		rows = append(rows, r)
	}
	return rows, nil
}

func (*GoogleSheetPortfolio) parseRow(ctx context.Context, index int, rawRow []interface{}) (*domain.PortfolioRow, error) {
	row := &domain.PortfolioRow{
		Index: index,
	}
	e := reflect.ValueOf(row).Elem()
	n := e.NumField()

	// Validate field lens: here we inject an index, and also leave the target optionl
	// meaning we should have the same number of columns as fields.
	if n <= len(rawRow) {
		slog.Info(ctx, "Expecting %v columns; got %v", e.NumField(), len(rawRow))
		return nil, terrors.BadRequest("invalid-column-schema", "Failed to parse row; invalid schema", map[string]string{
			"number_of_columns":          strconv.Itoa(len(rawRow)),
			"expected_number_of_columns": strconv.Itoa(e.NumField() - 1),
		})
	}

	for i, v := range rawRow {
		// We must convert the interface type here to string; this makes our life easier when marshalling
		// into the row struct.
		s, ok := v.(string)
		if !ok {
			return nil, terrors.BadRequest("bad-type", "Failed to parse portfolio row", nil)
		}

		if s == "" {
			continue
		}

		errParams := map[string]string{
			"column_index": strconv.Itoa(i),
		}

		// We want to the i-th field here, since we inject the row index into the row
		field := e.Field(i + 1)
		switch field.Kind() {
		case reflect.String:
			field.SetString(s)
		case reflect.Float64:
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, terrors.Augment(err, "Failed to parse row", errParams)
			}
			field.SetFloat(f)
		case reflect.Int:
			ii, err := strconv.ParseInt(s, 0, 64)
			if err != nil {
				return nil, terrors.Augment(err, "Failed to parse row", errParams)
			}
			field.SetInt(ii)
		default:
			return nil, terrors.BadRequest("invalid-type", "Cannot parse invalid type", errParams)
		}
	}
	return row, nil
}

func (*GoogleSheetPortfolio) parseMetadata(ctx context.Context, rawMetadata [][]interface{}) (*domain.PortfolioMetadata, error) {
	metadata := &domain.PortfolioMetadata{}
	e := reflect.ValueOf(metadata).Elem()

	// Validate metadata items (rows)
	// TODO: recursive check, each "row" should only contain the single item.
	if e.NumField() != len(rawMetadata) {
		slog.Info(ctx, "Expecting %v rows; got %v", e.NumField(), len(rawMetadata))
		return nil, terrors.BadRequest("invalid-metadata-schema", "Failed to parse metadata; invalid schema", map[string]string{
			"number_of_rows":          strconv.Itoa(len(rawMetadata)),
			"expected_number_of_rows": strconv.Itoa(e.NumField() - 1),
		})
	}
	// TODO refactor into own func
	for i, v := range rawMetadata {
		// We must convert the interface type here to string; this makes our life easier when marshalling
		// into the row struct.
		v0 := v[0]
		s, ok := v0.(string)
		if !ok {
			return nil, terrors.BadRequest("bad-type", "Failed to parse portfolio row", nil)
		}

		errParams := map[string]string{
			"row_number": strconv.Itoa(i),
		}

		field := e.Field(i)
		switch field.Kind() {
		case reflect.String:
			field.SetString(s)
		case reflect.Float64:
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, terrors.Augment(err, "Failed to parse metadata", errParams)
			}
			field.SetFloat(f)
		case reflect.Int:
			ii, err := strconv.ParseInt(s, 0, 64)
			if err != nil {
				return nil, terrors.Augment(err, "Failed to parse metadata", errParams)
			}
			field.SetInt(ii)
		default:
			return nil, terrors.BadRequest("invalid-type", "Cannot parse invalid type", errParams)
		}
	}
	return metadata, nil
}

func (GoogleSheetPortfolio) parseTradeHistoryRow(ctx context.Context, index int, values []interface{}) (*domain.HistoricalTradeRow, error) {
	row := &domain.HistoricalTradeRow{
		Index: index,
	}
	e := reflect.ValueOf(row).Elem()
	n := e.NumField()

	// Validate field lens: here we inject an index, and also leave the target optionl
	// meaning we should have the same number of columns as fields.
	if n <= len(values) {
		slog.Info(ctx, "Expecting %v columns; got %v", e.NumField(), len(values))
		return nil, terrors.BadRequest("invalid-column-schema", "Failed to parse trade history row; invalid schema", map[string]string{
			"number_of_columns":          strconv.Itoa(len(values)),
			"expected_number_of_columns": strconv.Itoa(e.NumField() - 1),
		})
	}

	for i, v := range values {
		// We must convert the interface type here to string; this makes our life easier when marshalling
		// into the row struct.
		s, ok := v.(string)
		if !ok {
			return nil, terrors.BadRequest("bad-type", "Failed to parse portfolio row", nil)
		}

		if s == "" {
			continue
		}

		errParams := map[string]string{
			"column_index": strconv.Itoa(i),
		}

		// We want to the i-th field here, since we inject the row index into the row
		field := e.Field(i + 1)
		switch field.Kind() {
		case reflect.String:
			field.SetString(s)
		case reflect.Float64:
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, terrors.Augment(err, "Failed to parse row", errParams)
			}
			field.SetFloat(f)
		case reflect.Int:
			ii, err := strconv.ParseInt(s, 0, 64)
			if err != nil {
				return nil, terrors.Augment(err, "Failed to parse row", errParams)
			}
			field.SetInt(ii)
		default:
			return nil, terrors.BadRequest("invalid-type", "Cannot parse invalid type", errParams)
		}
	}
	return row, nil
}

func (gsp *GoogleSheetPortfolio) UpdateRow(ctx context.Context, sheetID string, row *domain.PortfolioRow) error {
	// TODO move sheetID and rangeRow to the spreadsheet itself.
	rows := []*domain.PortfolioRow{row}
	return gsp.UpdateRows(ctx, sheetID, rows)
}

func (gsp *GoogleSheetPortfolio) UpdateRows(ctx context.Context, sheetID string, rows []*domain.PortfolioRow) error {
	r := formatRowsRange(sheetID)
	slog.Info(ctx, "Updating googlesheets range for portfolio rows: %s", r)

	values := [][]interface{}{}
	for _, row := range rows {
		values = append(values, row.ToArray())
	}
	err := gsp.c.UpdateRows(ctx, gsp.ID(), r, values)
	if err != nil {
		return terrors.Augment(err, "Failed to update portfolio rows", nil)
	}
	return nil
}

func (gsp *GoogleSheetPortfolio) UpdateMetadata(ctx context.Context, sheetID string, metadata *domain.PortfolioMetadata) error {
	r := formatMetadataRange(sheetID)
	slog.Info(ctx, "Updating googlesheets range for metadata: %s", r)
	values := [][]interface{}{
		{
			metadata.TotalPNL,
		},
		{
			metadata.TotalWorth,
		},
	}
	err := gsp.c.UpdateRows(ctx, gsp.ID(), r, values)
	if err != nil {
		return terrors.Augment(err, "Failed to update metadata", nil)
	}
	return nil
}

func (gsp *GoogleSheetPortfolio) UpdateTradeHistory(ctx context.Context, sheetID string, rows []*domain.HistoricalTradeRow) error {
	r := formatTradeHistoryRange(sheetID)
	slog.Info(ctx, "Updating googlesheets range for trade history: %s", r)

	values := [][]interface{}{}
	for _, row := range rows {
		values = append(values, row.ToArray())
	}
	return nil
}

func (gsp *GoogleSheetPortfolio) validateRows(rows []*domain.PortfolioRow) ([]*domain.PortfolioRow, error) {
	return rows, nil
}

func formatRowsRange(sheetID string) string {
	return fmt.Sprintf("%s!%s", sheetID, defaultPortfolioRowRange)
}

func formatMetadataRange(sheetID string) string {
	return fmt.Sprintf("%s!%s", sheetID, defaultPortfolioMetatdataRange)
}

func formatTradeHistoryRange(sheetID string) string {
	return fmt.Sprintf("%s!%s", sheetID, defaultPortfolioTradeHistoryRange)
}
