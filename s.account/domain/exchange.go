package domain

import "time"

// Exchange holds metadata for a given exchange.
type Exchange struct {
	ExchangeID string    `db:"exchange_id"`
	Exchange   string    `db:"exachange"`
	APIKey     string    `db:"api_key"`
	SecretKey  string    `db:"secret_key"`
	UserID     string    `db:"user_id"`
	Created    time.Time `db:"created"`
	Updated    time.Time `db:"updated"`
	IsActive   bool      `db:"is_active"`
}
