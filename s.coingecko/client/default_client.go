package client

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	coingecko "github.com/superoo7/go-gecko/v3"
)

type coingeckoClient struct {
	cli     *coingecko.Client
	coins   map[string]string
	coinsMu sync.RWMutex
}

// CoingeckoListCoinItem ...
type CoingeckoListCoinItem struct {
	Name   string
	Symbol string
	ID     string
}

func (c *coingeckoClient) GetAllCoinIDs(ctx context.Context) ([]*CoingeckoListCoinItem, error) {
	l, err := c.cli.CoinsList()
	if err != nil {
		return nil, terrors.Augment(err, "Failed to retreive coins list", nil)
	}

	coins := []*CoingeckoListCoinItem{}
	for _, coin := range *l {
		coins = append(coins, &CoingeckoListCoinItem{
			ID:     strings.ToLower(coin.ID),
			Name:   strings.ToLower(coin.Name),
			Symbol: strings.ToLower(coin.Symbol),
		})
	}

	return coins, nil
}

func (c *coingeckoClient) GetCoinInfoByID(ctx context.Context, id string) (*CoinRecord, error) {
	coin, err := c.cli.CoinsID(strings.ToLower(id), true, false, true, false, false, false)
	if err != nil {
		return nil, err
	}

	return &CoinRecord{
		LatestPrice:              coin.MarketData.CurrentPrice,
		PriceChangePercentage24h: coin.MarketData.PriceChangePercentage24hInCurrency,
		ATH:                      coin.MarketData.ATH,
	}, nil
}

func (c *coingeckoClient) Ping(ctx context.Context) error {
	if _, err := c.cli.Ping(); err != nil {
		return terrors.Augment(err, "Failed to establish connection to coingecko", nil)
	}

	return nil
}

func (c *coingeckoClient) RefreshCoins(ctx context.Context) {
	// Refresh loop that will get called every 24 hours; except the initial iteration.
	t := time.NewTicker(100 * time.Millisecond)
	var isFirstRefresh = true
	for {
		select {
		case <-t.C:
			var (
				coins    []*CoingeckoListCoinItem
				multiErr error
			)
			// Basic retry loop.
			for i := 0; i <= 5; i++ {
				cs, err := c.GetAllCoinIDs(ctx)
				if err != nil {
					multiErr = multierror.Append(multiErr, err)
					// Sleep incase we are rate limiting.
					time.Sleep(30 * time.Second)
					continue
				}
				coins = cs
				break
			}

			if len(coins) == 0 && multiErr != nil {
				slog.Error(ctx, "Failed after 5 retries to retrieve coingecko coin ids: errors %v ", multiErr)
			}

			c.coinsMu.Lock()
			for _, coin := range coins {
				if _, ok := blacklist[strings.ToLower(coin.ID)]; ok {
					slog.Info(ctx, "Skipping blacklisted coin: %s", coin.ID)
					continue
				}
				c.coins[coin.Symbol] = coin.ID
			}
			c.coinsMu.Unlock()

			if len(c.coins) == 0 {
				// We've retried 5 times; and the internal coins list is still empty, this means on start up we failed
				// to retrieve our list of coin id's. This service doesn't work without them. We should panic.
				panic("Failed to retreive set of coin id's from coingecko")
			}
		case <-ctx.Done():
			slog.Info(ctx, "Coingecko refresh token context cancelled: %v", ctx.Err())
			return
		}

		if isFirstRefresh {
			isFirstRefresh = false
			t.Reset(24 * time.Hour)
		}
	}
}

func (c *coingeckoClient) GetIDFromSymbol(symbol string) (string, error) {
	c.coinsMu.RLock()
	defer c.coinsMu.RUnlock()

	id, ok := c.coins[strings.ToLower(symbol)]
	if !ok {
		return "", terrors.BadResponse("failed-to-convert-symbol-to-id", "No id found for this symbol", map[string]string{
			"symbol": symbol,
		})
	}

	return id, nil
}
