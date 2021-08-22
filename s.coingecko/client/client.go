package client

import (
	"context"
	"net/http"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/ttlcache"
	"swallowtail/s.coingecko/cache"
	"time"

	"github.com/opentracing/opentracing-go"
	coingecko "github.com/superoo7/go-gecko/v3"
)

var (
	client CoinGeckoClient
)

// CoinGeckoClient ...
type CoinGeckoClient interface {
	// Ping checks the connectivity of the client, returns bool if we can reach the coingecko API.
	Ping(ctx context.Context) error
	// GetAllCoinIDs retrieves a list of all the coingecko coins; includes the symbol (ticker) & the coingecko ID for
	// reverse lookup.
	GetAllCoinIDs(ctx context.Context) ([]*CoingeckoListCoinItem, error)
	// GetCurrentPriceFromSymbol accepts a coin symbol (ticker) & returns the current price either from coingecko,
	// or the internal cache, if the value for that symbol hasn't expired. Also accepts an asset pair e.g `USDT`.
	GetCurrentPriceFromSymbol(ctx context.Context, symbol, assetPair string) (float64, error)
	// GetCurrentPriceFromID the same as GetCurrentPriceFromSymbol; however accepts the coingecko ID instead.
	// Also accepts an asset pair e.g `USDT`.
	GetCurrentPriceFromID(ctx context.Context, id, assetPair string) (float64, error)
	// GetATHFromSymbol retrieves the current ATH value from coingecko for the passed symbol (ticker).
	GetATHFromSymbol(ctx context.Context, symbol string) (float64, error)
	// GetATHFromID retrieves the current ATH value from coingecko for the passed coingecko id.
	GetATHFromID(ctx context.Context, id string) (float64, error)
	// RefreshCoins refreshes the internal cache of coin ids from coingecko.
	RefreshCoins(ctx context.Context)
}

// Init initializes the coingecko client.
func Init(ctx context.Context) error {
	// Create new cache.
	ttl := ttlcache.New(30 * time.Second)

	// Create new http client.
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	c := &coingeckoClient{
		cli: coingecko.NewClient(httpClient),
		cache: &cache.InMemoryCache{
			ttl: ttl,
		},
		coins: map[string]string{},
	}

	// Check connection is established.
	if err := c.Ping(ctx); err != nil {
		return gerrors.Augment(err, "failed_to_establish_coingecko_connection", nil)
	}

	// Kick off the background refresh loop.
	go c.RefreshCoins(ctx)

	client = c
	return nil
}

// GetCurrentPriceFromSymbol ...
func GetCurrentPriceFromSymbol(ctx context.Context, symbol, assetPair string) (float64, error) {
	return client.GetCurrentPriceFromSymbol(ctx, symbol, assetPair)
}

// GetCurrentPriceFromID ...
func GetCurrentPriceFromID(ctx context.Context, id, assetPair string) (float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get current price from coingecko")
	defer span.Finish()
	return client.GetCurrentPriceFromID(ctx, id, assetPair)
}

// GetATHFromSymbol ...
func GetATHFromSymbol(ctx context.Context, symbol string) (float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get ath price from coingecko")
	defer span.Finish()
	return client.GetATHFromSymbol(ctx, symbol)
}

// GetATHFromID ...
func GetATHFromID(ctx context.Context, id string) (float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get ath price from coingecko")
	defer span.Finish()
	return client.GetATHFromID(ctx, id)
}
