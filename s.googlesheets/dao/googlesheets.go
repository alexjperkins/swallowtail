package dao

import (
	"context"
	"swallowtail/s.googlesheets/domain"
	"swallowtail/s.googlesheets/templates"
	"time"

	"github.com/monzo/terrors"
)

// Register registers a googlesheet to a given user.
func RegisterGooglesheet(ctx context.Context, gs *domain.Googlesheet) error {
	var (
		sql = `
		INSERT INTO s_googlesheets_sheet
		(spreadsheet_id, sheet_id, user_id, with_pager_on_error, with_pager_on_target, created, updated)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
		`
	)

	now := time.Now().UTC()
	if _, err := (db.Exec(
		ctx, sql,
		gs.SpreadsheetID, gs.SheetID, gs.UserID, gs.WithPagerOnError, gs.WithPagerOnTarget,
		now, now,
	)); err != nil {
		return terrors.Propagate(err)
	}

	return nil
}

// ListSheetsByType retrieves all googlesheets registered of a particular type.
func ListSheetsByType(ctx context.Context, sheetType templates.SheetType) ([]*domain.Googlesheet, error) {
	var (
		sql = `
		SELECT * FROM s_googlesheets_sheet
		WHERE 
			googlesheet_type=$1
		`
		sheets = []*domain.Googlesheet{}
	)

	if err := db.Select(ctx, sheets, sql, sheetType.String()); err != nil {
		return nil, terrors.Propagate(err)
	}

	switch len(sheets) {
	case 0:
		return nil, terrors.NotFound("no-googlesheets-registered-with-this-type", "No googlesheets found with this type", nil)
	default:
		return sheets, nil
	}
}

// ListSheetsByUserID retrieves all sheets from persistent storage owned by the user id.
func ListSheetsByUserID(ctx context.Context, userID string) ([]*domain.Googlesheet, error) {
	var (
		sql = `
		SELECT * FROM s_googlesheets_sheet
		WHERE 
			user_id=$1
		`
		sheets = []*domain.Googlesheet{}
	)

	if err := db.Select(ctx, sheets, sql, userID); err != nil {
		return nil, terrors.Propagate(err)
	}

	switch len(sheets) {
	case 0:
		return nil, terrors.NotFound("no-googlesheets-registered-for-user", "No googlesheets found for user", nil)
	default:
		return sheets, nil
	}
}
