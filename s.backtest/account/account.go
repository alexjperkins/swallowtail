package domain

import (
	"fmt"
	"swallowtail/s.backtest/domain"
	"swallowtail/s.backtest/orderbook"
	"sync"
)

// New creates an account ready for backtesting
func New(accountSize float64, exchangeType string) *Account {
	var ob orderbook.Orderbook

	switch exchangeType {
	case orderbook.TypeDecentralizedExchange:
		ob = orderbook.NewAMMOrderbook()
	case orderbook.TypeCentralizedExchange:
		ob = orderbook.NewCEXOrderbook()
	default:
		panic(fmt.Sprintf("No orderbook for exchange type: %v", exchangeType))
	}

	return &Account{
		AccountSize:  accountSize,
		Balance:      accountSize,
		TradeHistory: TradeHistory{},
		orderbook:    ob,
	}
}

// Account
type Account struct {
	AccountSize     float64
	Balance         float64
	TradeHistory    TradeHistory
	tradeHistoryMtx sync.RWMutex
	orderbook       orderbook.Orderbook
}

// Add trdae adds trade to an account
func (a *Account) AddTrade(trade *domain.Trade) {
	a.tradeHistoryMtx.Lock()
	defer a.tradeHistoryMtx.Unlock()
	a.TradeHistory = append(a.TradeHistory, trade)
}

// Risk
func (a *Account) Risk(riskPercentage float64) float64 {
	return a.AccountSize * riskPercentage
}
