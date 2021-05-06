package orderbook

func NewCEXOrderbook() Orderbook {
	return &cexOrderbook{}
}

type cexOrderbook struct {
}

func (co *cexOrderbook) MakeTrade(trade *Trade) error {
	return nil
}
