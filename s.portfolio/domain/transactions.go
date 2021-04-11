package domain

import "time"

var (
	TransactionBuySide  = "BUY"
	TransactionSellSide = "SELL"
	TranactionsNull     = "NULL"
)

type Transaction struct {
	Asset   *Asset
	Side    string
	created time.Time
}

func NewTransaction(asset *Asset, side string) *Transaction {
	return &Transaction{
		Asset:   asset,
		Side:    side,
		created: time.Now(),
	}
}
