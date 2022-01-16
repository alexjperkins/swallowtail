package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/dao"
	"swallowtail/s.trade-engine/marshaling"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// ReadTradeStrategyByTradeStrategyID ...
func (s *TradeEngineService) ReadTradeStrategyByTradeStrategyID(
	ctx context.Context, in *tradeengineproto.ReadTradeStrategyByTradeStrategyIDRequest,
) (*tradeengineproto.ReadTradeStrategyByTradeStrategyIDResponse, error) {
	switch {
	case in.TradeStrategyId == "":
		return nil, gerrors.BadParam("missing_param.trade_strategy_id", nil)
	}

	errParams := map[string]string{
		"trade_strategy_id": in.TradeStrategyId,
	}

	trade, err := dao.ReadTradeStrategyByTradeStrategyID(ctx, in.TradeStrategyId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_trade_strategy_by_trade_strategy_id", errParams)
	}

	protoTrade := marshaling.TradeStrategyDomainToProto(trade)

	return &tradeengineproto.ReadTradeStrategyByTradeStrategyIDResponse{
		TradeStrategy: protoTrade,
	}, nil
}
