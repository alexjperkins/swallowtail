package domain

import "time"

// Trade ...
type Trade struct {
	ID             string
	ActorID        string
	ActorType      string
	IdempotencyKey string
	Exchange       string
	TradeType      string
	Asset          string
	Pair           string
	Entry          float64
	StopLoss       float64
	TakeProfits    []float64
	Status         string
	RiskReturn     float64
	Created        time.Time
	LastUpdate     time.Time
}

// TradeParticipent ...
type TradeParticipent struct {
	TradeID           string
	UserID            string
	Size              float64
	ExecutedTimestamp time.Time
}
