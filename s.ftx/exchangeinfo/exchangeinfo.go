package exchangeinfo

import (
	"context"
	"strings"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.ftx/client"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/monzo/slog"
)

var (
	mu          sync.RWMutex
	instruments = make(map[string]*client.Instrument)
)

// Init initializes the exchange info loop.
func Init(ctx context.Context) error {
	if err := gatherExchangeInfo(ctx); err != nil {
		return gerrors.Augment(err, "failed_to_init_exchange_info", nil)
	}

	go refresh(ctx)

	return nil
}

// GetInstrumentBySymbol ...
func GetInstrumentBySymbol(symbol string) (*client.Instrument, bool) {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := instruments[symbol]
	if !ok {
		return nil, false
	}

	return v, ok
}

func gatherExchangeInfo(ctx context.Context) error {
	var (
		rsp *client.ListInstrumentsResponse
		err error
	)
	for i := 0; i < 3; i++ {
		r, e := client.ListInstruments(ctx, &client.ListInstrumentsRequest{}, false)
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

	for _, instrument := range rsp.Instruments {
		instruments[strings.ReplaceAll(instrument.Symbol, "/", "")] = instrument
	}

	return nil
}

func refresh(ctx context.Context) {
	t := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-t.C:
			if err := gatherExchangeInfo(ctx); err != nil {
				slog.Error(ctx, "Failed to refresh binance exchange info: Error: %v", err)
				continue
			}
			slog.Info(ctx, "Refreshed binance info")
		case <-ctx.Done():
			return
		}
	}
}
