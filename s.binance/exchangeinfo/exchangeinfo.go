package exchangeinfo

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/monzo/slog"
)

var (
	quantityPrecisions = map[string]int{}
	pricePrecisions    = map[string]int{}
	mu                 sync.RWMutex
)

// Init ...
func Init(ctx context.Context) error {
	if err := gatherExchangeInfo(ctx); err != nil {
		return err
	}

	slog.Info(ctx, "Gathered required futures exchange information: %v %v", quantityPrecisions, pricePrecisions)

	// Start our refresh loop.
	go refresh(ctx)

	return nil
}

func refresh(ctx context.Context) {
	t := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-t.C:
			if err := gatherExchangeInfo(ctx); err != nil {
				slog.Error(ctx, "Failed to refresh exchange info: Error: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func gatherExchangeInfo(ctx context.Context) error {
	var (
		rsp *client.GetFuturesExchangeInfoResponse
		err error
	)
	for i := 0; i < 3; i++ {
		r, e := client.GetFuturesExchangeInfo(ctx, &client.GetFuturesExchangeInfoRequest{})
		if e != nil {
			multierror.Append(err, e)
			slog.Trace(ctx, "Failed to gather exchangeinfo, attempt [%v]; retrying...", i)
		}

		rsp = r
		break
	}

	if err != nil {
		return gerrors.Augment(err, "failed_to_init_exchange_info", nil)
	}

	if rsp == nil {
		return gerrors.Augment(err, "failed_to_init_exchange_info.empty_response", nil)
	}

	mu.Lock()
	defer mu.Unlock()

	for _, s := range rsp.Symbols {
		quantityPrecisions[s.BaseAsset] = s.QuantityPrecision
		pricePrecisions[s.BaseAsset] = s.PricePrecision
	}

	return nil
}

// GetBaseAssetQuantityPrecision ...
func GetBaseAssetQuantityPrecision(baseAsset string) (int, bool) {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := quantityPrecisions[baseAsset]
	if !ok {
		return 0, false
	}

	return v, true
}

// GetBaseAssetPricePrecision ...
func GetBaseAssetPricePrecision(baseAsset string) (int, bool) {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := pricePrecisions[baseAsset]
	if !ok {
		return 0, false
	}

	return v, true
}
