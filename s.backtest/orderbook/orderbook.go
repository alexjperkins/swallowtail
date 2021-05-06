package orderbook

var (
	TypeDecentralizedExchange = "decentralized-exchange"
	TypeCentralizedExchange   = "centralized-exchange"
)

type Orderbook interface {
	MakeTrade(trade *Trade) error
}
