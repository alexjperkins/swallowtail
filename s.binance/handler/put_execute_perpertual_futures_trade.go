package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	binanceproto "swallowtail/s.binance/proto"

	"github.com/monzo/slog"
)

// ExecuteFuturesPerpetualsTrade ...
func (s *BinanceService) ExecuteFuturesPerpetualsTrade(
	ctx context.Context, in *binanceproto.ExecuteFuturesPerpetualsTradeRequest,
) (*binanceproto.ExecuteFuturesPerpetualsTradeResponse, error) {
	switch {
	case in.Asset == "":
		return nil, gerrors.BadParam("missing_param.asset", nil)
	case in.Pair == "":
		return nil, gerrors.BadParam("missing_param.pair", nil)
	}

	slog.Warn(ctx, "%+v", in)

	return &binanceproto.ExecuteFuturesPerpetualsTradeResponse{
		ExchangeTradeId: "success",
	}, nil
}
