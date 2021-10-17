package client

import (
	"context"
	"net/http"
	"time"

	"github.com/monzo/slog"
	"github.com/opentracing/opentracing-go"
	coingecko "github.com/superoo7/go-gecko/v3"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.coingecko/cache"
)

var (
	client      CoinGeckoClient
	ttlcache    cache.CoingeckoCache
	rateLimiter *RateLimiter
)

// CoinInfo ...
type CoinRecord struct {
	LatestPrice              map[string]float64
	PriceChangePercentage24h map[string]float64
	ATH                      map[string]float64
}

// CoinGeckoClient ...
type CoinGeckoClient interface {
	// Ping checks the connectivity of the client, returns bool if we can reach the coingecko API.
	Ping(ctx context.Context) error
	// GetAllCoinIDs retrieves a list of all the coingecko coins; includes the symbol (ticker) & the coingecko ID for
	// reverse lookup.
	GetAllCoinIDs(ctx context.Context) ([]*CoingeckoListCoinItem, error)
	// Get GetCoinInfoByID fetches the latest coin info by id.
	GetCoinInfoByID(ctx context.Context, coinID string) (*CoinRecord, error)
	// Gets the Coingecko ID from a symbol passed.
	GetIDFromSymbol(symbol string) (string, error)
}

// Init initializes the coingecko client.
func Init(ctx context.Context) error {
	// Create new http client.
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Initialize client.
	c := &coingeckoClient{
		cli:   coingecko.NewClient(httpClient),
		coins: map[string]string{},
	}

	// Check connection is established & set.
	if err := c.Ping(ctx); err != nil {
		return gerrors.Augment(err, "failed_to_establish_coingecko_connection", nil)
	}
	client = c

	// Initialize rate limiter.
	rateLimiter = NewRateLimiter(ctx)

	// Initialize cache.
	ttlcache = cache.NewInMemoryCache(5*time.Minute, func(key string) (interface{}, time.Duration, error) {
		v, err := client.GetCoinInfoByID(ctx, key)
		if err != nil {
			return nil, 0, gerrors.Augment(err, "failed_to_fetch_latest_coin_info", map[string]string{
				"key": key,
			})
		}

		return v, 5 * time.Minute, nil
	})

	// Kick off the background refresh loop.
	go c.RefreshCoins(ctx)

	return nil
}

// GetCurrentPriceFromSymbol ...
func GetCurrentPriceFromSymbol(ctx context.Context, symbol, assetPair string) (float64, float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get current price from coingecko by symbol")
	defer span.Finish()

	// Thottle on rate limiter.
	rateLimiter.Throttle()

	errParams := map[string]string{
		"symbol":     symbol,
		"asset_pair": assetPair,
	}

	coinID, err := client.GetIDFromSymbol(symbol)
	if err != nil {
		return 0, 0, gerrors.Augment(err, "failed_to_get_current_price_from_symbol", errParams)
	}

	v, err := ttlcache.Get(coinID)
	if err != nil {
		return 0, 0, gerrors.Augment(err, "failed_to_get_current_price_from_symbol", errParams)
	}

	record, ok := v.(*CoinRecord)
	if !ok {
		slog.Warn(ctx, "Bad type coin gecko cache; failed to type assert", errParams)
		return 0, 0, gerrors.FailedPrecondition("failed_to_get_current_price_from_symbol.bad_record", errParams)
	}

	latestPrice, ok := record.LatestPrice[assetPair]
	if !ok {
		return 0, 0, gerrors.BadParam("failed_to_get_current_price_from_symbol.bad_asset_pair.latest_price", errParams)
	}

	percentagePriceChange24h, ok := record.PriceChangePercentage24h[assetPair]
	if !ok {
		return 0, 0, gerrors.BadParam("failed_to_get_current_price_from_symbol.bad_asset_pair.24h_change", errParams)
	}

	return latestPrice, percentagePriceChange24h, nil
}

// GetCurrentPriceFromID ...
func GetCurrentPriceFromID(ctx context.Context, coinID, assetPair string) (float64, float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get current price from coingecko by id")
	defer span.Finish()

	// Thottle on rate limiter.
	rateLimiter.Throttle()

	errParams := map[string]string{
		"coin_id":    coinID,
		"asset_pair": assetPair,
	}

	v, err := ttlcache.Get(coinID)
	if err != nil {
		return 0, 0, gerrors.Augment(err, "failed_to_get_current_price_from_id", errParams)
	}

	record, ok := v.(*CoinRecord)
	if !ok {
		slog.Warn(ctx, "Bad type coin gecko cache; failed to type assert", errParams)
		return 0, 0, gerrors.FailedPrecondition("failed_to_get_current_price_from_id.bad_record", errParams)
	}

	latestPrice, ok := record.LatestPrice[assetPair]
	if !ok {
		return 0, 0, gerrors.BadParam("failed_to_get_current_price_from_id.bad_asset_pair", errParams)
	}

	percentagePriceChange24h, ok := record.PriceChangePercentage24h[assetPair]
	if !ok {
		return 0, 0, gerrors.BadParam("failed_to_get_current_price_from_id.bad_asset_pair.24h_change", errParams)
	}

	return latestPrice, percentagePriceChange24h, nil
}

// GetATHFromSymbol ...
func GetATHFromSymbol(ctx context.Context, symbol, assetPair string) (float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get ath price from coingecko")
	defer span.Finish()

	// Thottle on rate limiter.
	rateLimiter.Throttle()

	errParams := map[string]string{
		"symbol":     symbol,
		"asset_pair": assetPair,
	}

	coinID, err := client.GetIDFromSymbol(symbol)
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_get_ath_from_symbol", errParams)
	}

	v, err := ttlcache.Get(coinID)
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_get_ath_from_symbol", errParams)
	}

	record, ok := v.(*CoinRecord)
	if !ok {
		slog.Warn(ctx, "Bad type coin gecko cache; failed to type assert", errParams)
		return 0, gerrors.FailedPrecondition("failed_to_get_ath_from_symbo.bad_record", errParams)
	}

	ath, ok := record.ATH[assetPair]
	if !ok {
		return 0, gerrors.BadParam("failed_to_get_ath_from_symbol.bad_asset_pair", errParams)
	}

	return ath, nil
}

// GetATHFromID ...
func GetATHFromID(ctx context.Context, coinID, assetPair string) (float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get ath price from coingecko")
	defer span.Finish()

	// Thottle on rate limiter.
	rateLimiter.Throttle()

	errParams := map[string]string{
		"coin_id":    coinID,
		"asset_pair": assetPair,
	}

	v, err := ttlcache.Get(coinID)
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_get_ath_from_id", errParams)
	}

	record, ok := v.(*CoinRecord)
	if !ok {
		slog.Warn(ctx, "Bad type coin gecko cache; failed to type assert", errParams)
		return 0, gerrors.FailedPrecondition("failed_to_get_ath_from_symbo.bad_record", errParams)
	}

	ath, ok := record.ATH[assetPair]
	if !ok {
		return 0, gerrors.BadParam("failed_to_get_ath_from_id.bad_asset_pair", errParams)
	}

	return ath, nil
}
