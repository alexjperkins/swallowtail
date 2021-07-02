package domain

import "time"

// Googlesheets holds metadata for accounts that use googlesheets.
type GooglesheetMetadata struct {
	ID            string    `db:"googlesheets_id"`
	SpreadsheetID string    `db:"spreadsheet_id"`
	SheetID       string    `db:"sheet_id"`
	Created       time.Time `db:"created"`
	Updated       time.Time `db:"updated"`
}
