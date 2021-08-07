package domain

import (
	"time"
)

// Googlesheet ...
type Googlesheet struct {
	GooglesheetID     string    `pb:"googlesheet_id"`
	SpreadsheetID     string    `pb:"spreadsheet_id"`
	SheetID           string    `pb:"sheet_id"`
	Email             string    `pb:"email"`
	UserID            string    `pb:"user_id"`
	SheetType         string    `pb:"sheet_type"`
	WithPagerOnError  bool      `pb:"with_pager_on_error"`
	WithPagerOnTarget bool      `pb:"with_pager_on_target"`
	Created           time.Time `pb:"created"`
	Updated           time.Time `pb:"updated"`
	Active            bool      `pb:"active"`
	URL               string    `pb:"url"`
}
