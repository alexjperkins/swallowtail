package handler

import (
	"context"
	"strconv"
	"swallowtail/libraries/gerrors"
	binanceclient "swallowtail/s.binance/client"
	"swallowtail/s.binance/domain"
	binanceproto "swallowtail/s.binance/proto"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/monzo/terrors"
)

// executeTradeWithRetry attempts to execute a trade with the given executor, in a retry loop.
// we exit the retry loop if:
//
// - max attempts reached.  - the deadline is reached to make an attempt.
// - we cannot handle the execution error
//
func executeTradeWithRetry(ctx context.Context, executer func(context.Context, *domain.Trade) error, trade *domain.Trade, maxAttempts int) (*domain.Trade, error) {
	tradeCtx, cancel := context.WithDeadline(ctx, trade.AttemptRetryUntil)
	defer cancel()

	var attempts int
	boff := backoff.NewExponentialBackOff()
	for {
		// Attempt to execute the trade 5 times.
		if attempts > maxAttempts {
			break
		}
		// Check the deadline first; if our trade is latent then we don't want to execute it
		// if it's already passed the deadline.
		select {
		case <-ctx.Done():
			// We didn't manage to execute the trade before the deadline.
			return nil, nil
		default:
			attempts++
		}

		// Attempt to make the trade.
		trade.Attempted = time.Now()
		err := binanceclient.ExecuteSpotTrade(tradeCtx, trade)
		switch {
		case terrors.Is(err, terrors.ErrRateLimited):
			// We've been rate limited; lets sleep based on an exponetial backoff algorithm.
			// TODO: does binance return rate limit data?
			time.Sleep(boff.NextBackOff())
		case err != nil:
			// We have an error that we can't handle.
			return nil, terrors.Augment(err, "Failed to execute trade; not retrying", map[string]string{
				"attempt":             strconv.Itoa(attempts),
				"attempted_timestamp": trade.Attempted.String(),
			})
		}
		// We executed the trade we can now exit the loop.
		break
	}

	return trade, nil
}

func isValidActor(actorID string) bool {
	switch actorID {
	case binanceproto.BinanceAccountActorManual, binanceproto.BinanceAccountActorTradeEngineSystem:
		return true
	default:
		return false
	}
}

func isValidCredentials(credentials *binanceproto.Credentials, apiKeyOnly bool) error {
	switch {
	case credentials == nil:
		return gerrors.BadParam("missing_param.credentials", nil)
	case credentials.ApiKey == "":
		return gerrors.BadParam("missing_param.credentials.api_key", nil)
	case !apiKeyOnly && credentials.SecretKey == "":
		return gerrors.BadParam("missing_param.credentials.secret_key", nil)
	default:
		return nil
	}
}
