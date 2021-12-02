package algo

import (
	"context"

	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func init() {
	register(tradeengineproto.EXECUTION_STRATEGY_DMA_MARKET, &DMAMarket{})
}

type DMAMarket struct{}

func (d *DMAMarket) Execute(ctx context.Context, strategy *tradeengineproto.TradeStrategy) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
	return nil, nil
}
