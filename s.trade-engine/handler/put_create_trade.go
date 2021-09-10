package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// CreateTrade ...
func CreateTrade(
	ctx context.Context, in *tradeengineproto.CreateTradeRequest,
) (*tradeengineproto.CreateTradeResponse, error) {
	switch {
	case in.Trade == nil:
		return nil, gerrors.BadParam("missing_param.trade", nil)
	}

	return nil, nil
}
