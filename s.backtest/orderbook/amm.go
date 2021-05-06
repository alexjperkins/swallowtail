package orderbook

// NewAMMOrderbook returns a new AMM orderbook
func NewAMMOrderbook() Orderbook {
	return &ammOrderbook{}
}

type ammOrderbook struct {
}

func (ao *ammOrderbook) MakeTrade(trade *Trade) error {
	return nil
}
