package execution

import (
	"context"

	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func init() {
	register(tradeengineproto.EXECUTION_STRATEGY_DMA_LIMIT, &DMALimit{})
}

// DMALimit ...
type DMALimit struct{}

// Execute ...
func (d *DMALimit) Execute(ctx context.Context, strategy *tradeengineproto.TradeStrategy, participant *tradeengineproto.ExecuteTradeStrategyForParticipantRequest) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
	return nil, nil
}
