package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/dao"
	"swallowtail/s.trade-engine/marshaling"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// ReadTradeByTradeID ...
func (s *TradeEngineService) ReadTradeByTradeID(
	ctx context.Context, in *tradeengineproto.ReadTradeByTradeIDRequest,
) (*tradeengineproto.ReadTradeByTradeIDResponse, error) {
	switch {
	case in.TradeId == "":
		return nil, gerrors.BadParam("missing_param.trade_id", nil)
	}

	errParams := map[string]string{
		"trade_id": in.TradeId,
	}

	trade, err := dao.ReadTradeByTradeID(ctx, in.TradeId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_trade_by_trade_id", errParams)
	}

	protoTrade := marshaling.TradeDomainToProto(trade)

	return &tradeengineproto.ReadTradeByTradeIDResponse{
		Trade: protoTrade,
	}, nil
}
