package execution

import (
	"context"

	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func init() {
	register(tradeengineproto.EXECUTION_STRATEGY_DMA_MARKET, &DMAMarket{})
}

// DMAMarket ...
type DMAMarket struct{}

// Execute ...
func (d *DMAMarket) Execute(ctx context.Context, strategy *tradeengineproto.TradeStrategy, participant *tradeengineproto.ExecuteTradeStrategyForParticipantRequest) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
	return nil, nil
}
