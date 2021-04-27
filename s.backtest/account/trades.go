package domain

import "time"

type Trade struct {
	Ticker      string
	Type        string
	Entry       float64
	StopLosses  []float64
	TakeProfits []float64
	PNL         float64
	TradeSize   float64
	Reason      string
	Opened      time.Time
	Closed      time.Time
}

type TradeHistory []*Trade

func (t TradeHistory) StrikeRate() float64 {
	var won int
	for _, trade := range t {
		if trade.PNL > 0 {
			won++
			continue
		}
		won--
	}
	return float64(won / len(t))
}
