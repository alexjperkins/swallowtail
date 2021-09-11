package handler

import (
	"context"
	"strconv"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	binanceproto "swallowtail/s.binance/proto"
)

// GetLatestPrice ...
func (s *BinanceService) GetLatestPrice(
	ctx context.Context, in *binanceproto.GetLatestPriceRequest,
) (*binanceproto.GetLatestPriceResponse, error) {
	switch {
	case in.Symbol == "":
		return nil, gerrors.BadParam("missing_param.symbol", nil)
	}

	errParams := map[string]string{
		"symbol": in.Symbol,
	}

	rsp, err := client.GetLatestPrice(ctx, &client.GetLatestPriceRequest{
		Symbol: in.Symbol,
	})

	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_latest_price.client_failure", errParams)
	}

	price, err := strconv.ParseFloat(rsp.Price, 32)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_latest_price.bad_price", errParams)
	}

	return &binanceproto.GetLatestPriceResponse{
		Price:     float32(price),
		Timestamp: int64(rsp.Time),
	}, nil
}
