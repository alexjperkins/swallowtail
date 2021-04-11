package indicators

import (
	"context"
	"strconv"
	"time"

	"swallowtail/s.binance-consumer/domain"

	"github.com/monzo/slog"
)

var (
	SSLIndicatorID = "ssl_channel_indicator"

	SSLClosePriceDimension = "close_price"
	SSLOpenPriceDimension  = "open_price"
	SSLDimensions          = []string{SSLClosePriceDimension, SSLOpenPriceDimension}

	SSLDownTrendSignal = "ssl_channel_crossover_downtrend"
	SSLUpTrendSignal   = "ssl_channel_crossover_uptrend"

	SSLInput = "ssl_input"
)

func SSLChannelIndicator(ctx context.Context, bufSize int, period int) *Indicator {
	return IndicatorFactory(ctx, SSLIndicatorID, bufSize, period, sslChannelSignalGenerator, SSLDimensions)
}

func sslChannelSignalGenerator(ctx context.Context, event interface{}, output chan *Signal, cache Cache) {
	// TODO: do we need a mutex here?
	previousCloseMean, err := cache[SSLClosePriceDimension].Mean()
	if err != nil {
		// Best effort
		slog.Error(ctx, "Failed to calcute close price mean: %v", err)
		return
	}
	previousOpenMean, err := cache[SSLOpenPriceDimension].Mean()
	if err != nil {
		// Best effort
		slog.Error(ctx, "Failed to calcute open price mean: %v", err)
		return
	}
	trendingUp := previousCloseMean > previousOpenMean

	switch e := event.(type) {
	case *domain.BinanceKlineEvent:
		// Calculate signal
		slog.Warn(ctx, "SSL event received -> %v", e)
		switch {
		case e.Data.ClosePrice <= e.Data.OpenPrice && trendingUp:
			// We have a cross over here into a downtrend.
			s := &Signal{
				Side:       NoSideSignal,
				Reason:     SSLDownTrendSignal,
				Timestamp:  time.Now(),
				Confidence: 1.0,
			}
			select {
			case output <- s:
				slog.Info(ctx, "Downtrend signal determined: %v", s)
			default:
				slog.Warn(ctx, "Signal channel busy; couldn't push signal: %v", s)
			}
		case e.Data.ClosePrice > e.Data.OpenPrice && !trendingUp:
			// We have a cross over here into an uptrend.
			s := &Signal{
				Side:       NoSideSignal,
				Reason:     SSLUpTrendSignal,
				Timestamp:  time.Now(),
				Confidence: 1.0,
			}
			select {
			case output <- s:
				slog.Info(ctx, "Uptrend signal determined: %v", s)
			default:
				slog.Warn(ctx, "Signal channel busy; couldn't push signal: %v", s)
			}
		default:
			// No signal found; continue about our day.
		}
		// Now that we've calculated our signal, we can push the new event into the cache
		// This is actually wrong. We need to recalculate the means here.
		cp, err := strconv.ParseFloat(e.Data.ClosePrice, 32)
		if err != nil {
			// Best effort
			slog.Error(ctx, "Failed to parse closing price: %v", e.Data.ClosePrice)
			return
		}
		op, err := strconv.ParseFloat(e.Data.ClosePrice, 32)
		if err != nil {
			// Best effort
			slog.Error(ctx, "Failed to parse opening price: %v", e.Data.ClosePrice)
			return
		}
		cache[SSLClosePriceDimension].Push(float32(cp))
		cache[SSLOpenPriceDimension].Push(float32(op))
	default:
		slog.Error(ctx, "Malformed event cannot parse: %v", event)
		return
	}
}
