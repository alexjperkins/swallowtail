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
)

// getLatestPrice ...
func getLatestPrice(ctx context.Context, asset, pair string) (float64, error) {
	var merr multierror.Error

	// First try binance to get the latest price. Here we use a retry loop if we do get rate limited.
	for i := 0; i < 3; i++ {
		rsp, err := (&binanceproto.GetLatestPriceRequest{
			Symbol: fmt.Sprintf("%s%s", asset, pair),
		}).SendWithTimeout(ctx, 30*time.Second).Response()

		boff := backoff.NewExponentialBackOff()
		switch {
		case gerrors.Is(err, gerrors.ErrRateLimited):
			d := boff.NextBackOff()
			slog.Trace(ctx, "Rate Limited: %v; sleeping for %s", err, d)
			time.Sleep(d)
		case err != nil:
			merr.Errors = append(merr.Errors, err)
			break
		}

		return float64(rsp.Price), nil
	}

	// Fallback option if we fail on binance; we then try coingecko.
	rsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
		AssetPair:   pair,
		AssetSymbol: asset,
	}).SendWithTimeout(ctx, 15*time.Second).Response()
	if err != nil {
		merr.Errors = append(merr.Errors, err)
		return 0, &merr
	}

	return float64(rsp.LatestPrice), nil
}
