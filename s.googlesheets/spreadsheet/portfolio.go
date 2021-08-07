package spreadsheet

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"swallowtail/s.googlesheets/client"
	"swallowtail/s.googlesheets/domain"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

const (
	defaultPortfolioRowRange          = "C8:L100"
	defaultPortfolioMetatdataRange    = "D2:D4"
	defaultPortfolioTradeHistoryRange = "Q8:W100"
)

// Portfolio ...
type Portfolio struct{}

func (gsp *Portfolio) Rows(ctx context.Context, spreadsheetID, sheetID string) ([]*domain.PortfolioRow, error) {
	r := formatRowsRange(sheetID)
	rows := []*domain.PortfolioRow{}

	vs, err := client.Values(ctx, spreadsheetID, r)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to retreive sheet values", map[string]string{
			"sheet_id":       sheetID,
			"spreadsheet_id": spreadsheetID,
			"range":          r,
		})
	}
	if vs == nil {
		return rows, nil
	}

	for i, v := range vs {
		if len(v) == 0 {
			continue
		}

		row, err := gsp.parseRow(ctx, i, v)
		if err != nil {
			return nil, terrors.Augment(err, "Failed to parse row", map[string]string{
				"sheet_id":       sheetID,
				"range":          r,
				"spreadsheet_id": spreadsheetID,
			})
		}
		rows = append(rows, row)
	}
	return gsp.validateRows(rows)
}

func (gsp *Portfolio) Metadata(ctx context.Context, spreadsheetID, sheetID string) (*domain.PortfolioMetadata, error) {
	r := formatMetadataRange(sheetID)

	vs, err := client.Values(ctx, spreadsheetID, r)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to retrieve sheet pnl metadata", map[string]string{
			"sheet_id":       sheetID,
			"range":          r,
			"spreadsheet_id": spreadsheetID,
		})
	}
	if vs == nil {
		return &domain.PortfolioMetadata{}, nil
	}

	return gsp.parseMetadata(ctx, vs)
}

func (gsp *Portfolio) TradeHistory(ctx context.Context, spreadsheetID, sheetID string) ([]*domain.HistoricalTradeRow, error) {
	r := formatTradeHistoryRange(sheetID)

	rows := []*domain.HistoricalTradeRow{}

	vs, err := client.Values(ctx, spreadsheetID, r)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to retrieve sheet trade history", map[string]string{
			"sheet_id":       sheetID,
			"range":          r,
			"spreadsheet_id": spreadsheetID,
		})

	}
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
				"sheet_id":       sheetID,
				"spreadsheet_id": spreadsheetID,
			})
		}
		rows = append(rows, r)
	}
	return rows, nil
}

func (*Portfolio) parseRow(ctx context.Context, index int, rawRow []interface{}) (*domain.PortfolioRow, error) {
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

func (*Portfolio) parseMetadata(ctx context.Context, rawMetadata [][]interface{}) (*domain.PortfolioMetadata, error) {
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
		var s string
		if len(v) == 1 {
			v0 := v[0]
			s, _ = v0.(string)
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

func (Portfolio) parseTradeHistoryRow(ctx context.Context, index int, values []interface{}) (*domain.HistoricalTradeRow, error) {
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

func (gsp *Portfolio) UpdateRow(ctx context.Context, spreadsheetID, sheetID string, row *domain.PortfolioRow) error {
	rows := []*domain.PortfolioRow{row}
	return gsp.UpdateRows(ctx, spreadsheetID, sheetID, rows)
}

func (gsp *Portfolio) UpdateRows(ctx context.Context, spreadsheetID, sheetID string, rows []*domain.PortfolioRow) error {
	r := formatRowsRange(sheetID)

	values := [][]interface{}{}
	for _, row := range rows {
		values = append(values, row.ToArray())
	}

	if err := client.UpdateRows(ctx, spreadsheetID, r, values); err != nil {
		return terrors.Augment(err, "Failed to update portfolio rows", nil)
	}
	return nil
}

func (gsp *Portfolio) UpdateMetadata(ctx context.Context, spreadsheetID, sheetID string, metadata *domain.PortfolioMetadata) error {
	r := formatMetadataRange(sheetID)

	values := [][]interface{}{
		{
			metadata.TotalPNL,
		},
		{
			metadata.TotalWorth,
		},
	}

	if err := client.UpdateRows(ctx, spreadsheetID, r, values); err != nil {
		return terrors.Augment(err, "Failed to update metadata", nil)
	}
	return nil
}

func (gsp *Portfolio) UpdateTradeHistory(ctx context.Context, spreadsheetID, sheetID string, rows []*domain.HistoricalTradeRow) error {
	r := formatTradeHistoryRange(sheetID)

	values := [][]interface{}{}
	for _, row := range rows {
		values = append(values, row.ToArray())
	}

	if err := client.UpdateRows(ctx, spreadsheetID, r, values); err != nil {
		return terrors.Augment(err, "Failed to update trade history", nil)
	}

	return nil
}

func (gsp *Portfolio) validateRows(rows []*domain.PortfolioRow) ([]*domain.PortfolioRow, error) {
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
