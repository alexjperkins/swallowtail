package domain

import (
	"sync"
)

// New creates an account ready for backtesting
func New(accountSize float64) *Account {
	return &Account{
		AccountSize:  accountSize,
		Balance:      accountSize,
		TradeHistory: TradeHistory{},
	}
}

// Account
type Account struct {
	AccountSize     float64
	Balance         float64
	TradeHistory    TradeHistory
	tradeHistoryMtx sync.RWMutex
}

// Add trdae adds trade to an account
func (a *Account) AddTrade(trade *Trade) {
	a.tradeHistoryMtx.Lock()
	defer a.tradeHistoryMtx.Unlock()
	a.TradeHistory = append(a.TradeHistory, trade)
}

// Risk
func (a *Account) Risk(riskPercentage float64) float64 {
	return a.AccountSize * riskPercentage
}
