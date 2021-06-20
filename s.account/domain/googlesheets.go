package domain

// Googlesheets holds metadata for accounts that use googlesheets.
type GooglesheetMetadata struct {
	ID            string `db:"googlesheets_id"`
	SpreadsheetID string `db:"spreadsheet_id"`
	SheetID       string `db:"sheet_id"`
}
