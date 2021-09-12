package domain

import "time"

// Trade ...
type Trade struct {
	TradeID            string
	ActorID            string
	HumanizedActorName string
	ActorType          string
	IdempotencyKey     string
	OrderType          string
	TradeType          string
	Asset              string
	Pair               string
	Entry              float64
	StopLoss           float64
	TakeProfits        []float64
	Status             string
	Created            time.Time
	LastUpdated        time.Time
	TradeSide          string
	CurrentPrice       float64
}

// TradeParticipent ...
type TradeParticipent struct {
	TradeID           string
	UserID            string
	IsBot             bool
	Size              float64
	Exchange          string
	ExchangeOrderID   string
	Status            string
	ExecutedTimestamp time.Time
}
