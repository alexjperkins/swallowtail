package domain

import "time"

// Trade ...
type Trade struct {
	TradeID            string    `db:"trade_id"`
	ActorID            string    `db:"actor_id"`
	HumanizedActorName string    `db:"humanized_actor_name"`
	ActorType          string    `db:"actor_type"`
	IdempotencyKey     string    `db:"idempotency_key"`
	OrderType          string    `db:"order_type"`
	TradeType          string    `db:"trade_type"`
	Asset              string    `db:"asset"`
	Pair               string    `db:"pair"`
	Entries            []float64 `db:"entries"`
	StopLoss           float64   `db:"stop_loss"`
	TakeProfits        []float64 `db:"take_profits"`
	Status             string    `db:"status"`
	Created            time.Time `db:"created"`
	LastUpdated        time.Time `db:"last_updated"`
	TradeSide          string    `db:"trade_side"`
	CurrentPrice       float64   `db:"current_price"`
	TradeableExchanges []string  `db:"tradeable_exchanges"`
}

// TradeParticipant ...
type TradeParticipant struct {
	TradeParticipantID string    `db:"trade_participant_id"`
	TradeID            string    `db:"trade_id"`
	UserID             string    `db:"user_id"`
	IsBot              bool      `db:"is_bot"`
	Size               float64   `db:"size"`
	Risk               float64   `db:"risk"`
	Exchange           string    `db:"exchange"`
	ExchangeOrderID    string    `db:"exchange_order_id"`
	Status             string    `db:"status"`
	ExecutedTimestamp  time.Time `db:"executed"`
	DCAStrategy        string    `db:"dca_strategy"`
}
