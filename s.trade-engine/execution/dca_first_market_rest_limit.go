package algo

import (
	"context"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func init() {
	register(tradeengineproto.EXECUTION_STRATEGY_DCA_FIRST_MARKET_REST_LIMIT, &DCAFirstMarketRestLimit{})
}

type DCAFirstMarketRestLimit struct{}

func (d *DCAFirstMarketRestLimit) Execute(ctx context.Context, strategy *tradeengineproto.TradeStrategy) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
	return nil, nil
}
