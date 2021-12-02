package domain

import "time"

// TradeStrategy ...
type TradeStrategy struct {
	TradeStrategyID    string    `db:"trade_strategy_id"`
	ActorID            string    `db:"actor_id"`
	HumanizedActorName string    `db:"humanized_actor_name"`
	ActorType          string    `db:"actor_type"`
	IdempotencyKey     string    `db:"idempotency_key"`
	ExecutionStrategy  string    `db:"execution_strategy"`
	InstrumentType     string    `db:"instrument_type"`
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
	TradeableVenues    []string  `db:"tradeable_venues"`
}

// TradeStrategyParticipant ...
type TradeStrategyParticipant struct {
	ParticipantID     string    `db:"participant_id"`
	TradeStrategyID   string    `db:"trade_strategy_id"`
	UserID            string    `db:"user_id"`
	IsBot             bool      `db:"is_bot"`
	Size              float64   `db:"size"`
	Risk              float64   `db:"risk"`
	Venue             string    `db:"venue"`
	ExchangeOrderIDs  []string  `db:"exchange_order_ids"`
	Status            string    `db:"status"`
	ExecutedTimestamp time.Time `db:"executed"`
}
