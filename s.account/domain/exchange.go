package domain

// Exchange holds metadata for a given exchange.
type Exchange struct {
	ID        string `db:"exchange_id"`
	Exchange  string `db:"exachange"`
	APIKey    string `db:"api_key"`
	SecretKey string `db:"secret_key"`
}
