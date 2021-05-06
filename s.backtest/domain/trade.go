package domain

import "time"

var (
	TypeLimitOrderBuySide   = "limit-order-buy-side"
	TypeLimitOrderSellSide  = "limit-order-sell-side"
	TypeMarketOrderBuySide  = "market-order-buy-side"
	TypeMarketOrderSellSide = "market-order-sell-side"
	TypeStopLossBuySide     = "stop-loss-buy-side"
	TypeStopLossSellSide    = "stop-loss-sell-side"
	TypeTakeProfitBuySide   = "take-profit-buy-side"
	TypeTakeProfitSellSide  = "take-profit-sell-side"

	StatusInProgress = "in-progress"
	StatusComplete   = "complete"
	StatusClosed     = "closed"
	StatusCancelled  = "cancelled"
)

// Trade models a trade on either centralized or decentralized exchanges.
type Trade struct {
	Type        string
	Pair        string
	Margin      float64
	Leverage    int
	Entry       float64
	StopLoss    *Trade
	TakeProfits []*Trade
	PNL         float64
	ReduceOnly  bool
	Status      string
	Fees        float64
	FundingFees float64
	Timestamp   time.Time
}
