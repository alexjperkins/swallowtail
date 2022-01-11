package domain

import (
	"time"
)

// VenueAccount holds metadata for a given exchange.
type VenueAccount struct {
	VenueAccountID string    `db:"venue_account_id"`
	VenueID        string    `db:"venue_id"`
	APIKey         string    `db:"api_key"`
	SecretKey      string    `db:"secret_key"`
	SubAccount     string    `db:"subaccount"`
	UserID         string    `db:"user_id"`
	Created        time.Time `db:"created"`
	Updated        time.Time `db:"updated"`
	IsActive       bool      `db:"is_active"`
	AccountAlias   string    `db:"account_alias"`
	URL            string    `db:"url"`
	WSURL          string    `db:"ws_url"`
}

// InternalVenueAccount holds metadata for internal venue accounts.
type InternalVenueAccount struct {
	VenueAccountID   string    `db:"venue_account_id"`
	VenueID          string    `db:"venue_id"`
	APIKey           string    `db:"api_key"`
	SecretKey        string    `db:"secret_key"`
	SubAccount       string    `db:"subaccount"`
	URL              string    `db:"url"`
	WSURL            string    `db:"ws_url"`
	VenueAccountType string    `db:"venue_account_type"`
	Created          time.Time `db:"created"`
	Updated          time.Time `db:"updated"`
}
