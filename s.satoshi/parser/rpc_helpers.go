package parser

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	binanceproto "swallowtail/s.binance/proto"
	coingeckoproto "swallowtail/s.coingecko/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// getLatestPrice ...
func getLatestPrice(ctx context.Context, asset string) (float64, error) {
	var merr multierror.Error

	// First try binance to get the latest price. Here we use a retry loop if we do get rate limited.
outer:
	for i := 0; i < 3; i++ {
		rsp, err := (&binanceproto.GetLatestPriceRequest{
			Symbol: fmt.Sprintf("%s%s", asset, tradeengineproto.TRADE_PAIR_USDT.String()),
		}).SendWithTimeout(ctx, 30*time.Second).Response()

		boff := backoff.NewExponentialBackOff()
		switch {
		case gerrors.Is(err, gerrors.ErrRateLimited):
			d := boff.NextBackOff()
			slog.Trace(ctx, "Rate Limited: %v; sleeping for %s", err, d)
			time.Sleep(d)
			continue
		case err != nil:
			merr.Errors = append(merr.Errors, err)
			break outer
		}

		return float64(rsp.Price), nil
	}

	slog.Warn(ctx, "Failed to fetch latest price for [%s] on Binance after 3 retries; attempting coingecko as fallback.", asset+tradeengineproto.TRADE_PAIR_USDT.String())

	// Fallback option if we fail on binance; we then try coingecko.
	rsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
		AssetPair:   tradeengineproto.TRADE_PAIR_USD.String(),
		AssetSymbol: asset,
	}).SendWithTimeout(ctx, 15*time.Second).Response()
	if err != nil {
		merr.Errors = append(merr.Errors, err)
		return 0, &merr
	}

	return float64(rsp.LatestPrice), nil
}
