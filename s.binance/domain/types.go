package domain

import "time"

// Trade is an abstraction for a given trade made to binance; it includes spot, perpetuals & quarterly futures
// TODO: move trade type to enum.
type Trade struct {
	TradeID           string    `db:"trade_id"`
	UserDiscordID     string    `db:"user_discord_id"`
	IdempotencyKey    string    `db:"idempotency_key"`
	Side              string    `db:"trade_side"`
	Type              string    `db:"trade_type"`
	AssetPair         string    `db:"asset_pair"`
	Amount            string    `db:"amount"`
	Value             string    `db:"value"`
	Created           time.Time `db:"created"`
	Attempted         time.Time `db:"attempted"`
	AttemptRetryUntil time.Time `db:"attempt_retry_until"`
}
