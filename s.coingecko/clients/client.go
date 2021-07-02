package clients

import (
	"context"
	"net/http"
	"strings"
	"swallowtail/libraries/ttlcache"
	"sync"
	"time"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	coingecko "github.com/superoo7/go-gecko/v3"
)

var (
	CoingeckoClientID = "coingecko-client-id"
	coingeckoClient   CoinGeckoClient

	defaultCoinGeckoClientTimeout = time.Second * 30

	defaultPriceTTL = time.Duration(30 * time.Second)

	symbolToCoingeckoID = map[string]string{}
	symbolMu            sync.RWMutex

	once sync.Once

	blacklist = map[string]bool{
		"universe-token": true,
	}
)

// TODO: fix to a singleton here.
type CoinGeckoClient interface {
	// Ping checks the connectivity of the client, returns bool if we can reach the coingecko API.
	Ping(ctx context.Context) bool
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
}

type CoingeckoListCoinItem struct {
	Name   string
	Symbol string
	ID     string
}

func init() {
	ctx := context.Background()
	c := New(ctx)
	coins, err := c.GetAllCoinIDs(ctx)
	if err != nil {
		panic(terrors.Augment(err, "Failed to retreive coingecko coins: err", nil))
	}

	symbolMu.Lock()
	defer symbolMu.Unlock()
	slog.Info(ctx, "Fetching all coins listed on coingecko & creating a mapping.")
	for _, coin := range coins {
		if _, ok := blacklist[strings.ToLower(coin.ID)]; ok {
			slog.Info(ctx, "Skipping blacklisted coin: %s", coin.ID)
			continue
		}
		symbolToCoingeckoID[coin.Symbol] = coin.ID
	}
	slog.Info(ctx, "Fetched all coins: total: %v", len(symbolToCoingeckoID))
}

func New(ctx context.Context) CoinGeckoClient {
	once.Do(func() {
		httpClient := &http.Client{
			Timeout: defaultCoinGeckoClientTimeout,
		}
		ttl := ttlcache.New(defaultPriceTTL)
		cgc := &coinGeckoClient{
			c:   coingecko.NewClient(httpClient),
			ttl: ttl,
		}
		if ok := cgc.Ping(ctx); !ok {
			slog.Error(context.TODO(), "Failed to connect coingecko client")
			panic("Failed to connect coingecko client")
		} else {
			slog.Info(context.TODO(), "Created coingecko client")
		}
		coingeckoClient = cgc
	})
	return coingeckoClient
}

type coinGeckoClient struct {
	c   *coingecko.Client
	ttl *ttlcache.TTLCache
}

func (cgc *coinGeckoClient) GetATHFromID(ctx context.Context, id string) (float64, error) {
	coinID, err := cgc.c.CoinsID(strings.ToLower(id), true, false, true, false, false, false)
	if err != nil {
		return 0.0, err
	}

	return coinID.MarketData.ATH["usd"], nil
}

func (cgc *coinGeckoClient) GetATHFromSymbol(ctx context.Context, symbol string) (float64, error) {
	id, err := getIDFromSymbol(symbol)
	if err != nil {
		return 0.0, err
	}
	return cgc.GetATHFromID(ctx, id)
}

func (cgc *coinGeckoClient) GetCurrentPriceFromID(ctx context.Context, id, assetPair string) (float64, error) {
	p, hasExpired := cgc.ttl.Get(id)
	if p == nil || hasExpired {
		ssp, err := cgc.c.SimpleSinglePrice(strings.ToLower(id), strings.ToLower(assetPair))
		if err != nil {
			slog.Error(ctx, "id: %v,  asset_pair: %v", id, assetPair)
			return 0.0, terrors.Augment(err, "Failed to retreive current price", map[string]string{
				"coingecko_id": id,
				"asset_pair":   assetPair,
			})
		}
		slog.Trace(ctx, "Updating coingecko price cache for [%s].", id)
		latestPrice := float64(ssp.MarketPrice)
		cgc.ttl.Set(id, latestPrice)
		return latestPrice, nil
	}

	slog.Trace(ctx, "Retrieving price for [%s] from coingecko cache.", id)
	pf, ok := p.(float64)
	if !ok {
		slog.Error(ctx, "Failed to parse price value into float %v -> %v", p, pf)
		return 0.0, terrors.BadResponse("invalid-price-type", "Failed to convert cached price to float", map[string]string{
			"id": id,
		})
	}
	return pf, nil
}

func (cgc *coinGeckoClient) GetCurrentPriceFromSymbol(ctx context.Context, symbol, assetPair string) (float64, error) {
	id, err := getIDFromSymbol(symbol)
	if err != nil {
		return 0.0, err
	}
	return cgc.GetCurrentPriceFromID(ctx, id, assetPair)
}

func (cgc *coinGeckoClient) GetAllCoinIDs(ctx context.Context) ([]*CoingeckoListCoinItem, error) {
	l, err := cgc.c.CoinsList()
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

func (cgc *coinGeckoClient) Ping(ctx context.Context) bool {
	_, err := cgc.c.Ping()
	if err != nil {
		slog.Error(ctx, "Failed to connect to coingecko: %v", err)
	}
	return err == nil
}

func getIDFromSymbol(symbol string) (string, error) {
	symbolMu.RLock()
	defer symbolMu.RUnlock()
	id, ok := symbolToCoingeckoID[strings.ToLower(symbol)]
	if !ok {
		return "", terrors.BadResponse("failed-to-convert-symbol-to-id", "No id found for this symbol", map[string]string{
			"symbol": symbol,
		})
	}
	slog.Info(nil, "Coingecko mapping: received: %s, converting -> %s", symbol, id)
	return id, nil
}
