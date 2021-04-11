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
	coingeckoClient   *CoinGeckoClient

	defaultCoinGeckoClientTimeout = time.Second * 10

	defaultPriceTTL = time.Duration(30 * time.Second)

	symbolToCoingeckoID = map[string]string{}
	symbolMu            sync.RWMutex

	once sync.Once

	blacklist = map[string]bool{
		"universe-token": true,
	}
)

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
		panic("Failed to retreive coingecko coins")
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

func New(ctx context.Context) *CoinGeckoClient {
	once.Do(func() {
		httpClient := &http.Client{
			Timeout: defaultCoinGeckoClientTimeout,
		}
		ttl := ttlcache.New(defaultPriceTTL)
		cgc := &CoinGeckoClient{
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

type CoinGeckoClient struct {
	c   *coingecko.Client
	ttl *ttlcache.TTLCache
}

func (cgc *CoinGeckoClient) GetATHFromID(ctx context.Context, id string) (float64, error) {
	coinID, err := cgc.c.CoinsID(strings.ToLower(id), true, false, true, false, false, false)
	if err != nil {
		return 0.0, err
	}

	return coinID.MarketData.ATH["usd"], nil
}

func (cgc *CoinGeckoClient) GetCurrentPriceFromID(ctx context.Context, id, assetPair string) (float64, error) {
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

func (cgc *CoinGeckoClient) GetCurrentPriceFromSymbol(ctx context.Context, symbol, assetPair string) (float64, error) {
	id, err := getIDFromSymbol(symbol)
	if err != nil {
		return 0.0, err
	}
	return cgc.GetCurrentPriceFromID(ctx, id, assetPair)
}

func (cgc *CoinGeckoClient) GetAllCoinIDs(ctx context.Context) ([]*CoingeckoListCoinItem, error) {
	l, err := cgc.c.CoinsList()
	if err != nil {
		return nil, terrors.Augment(err, "Failed to retreive coins list", nil)
	}
	coins := []*CoingeckoListCoinItem{}
	for _, coin := range *l {
		coins = append(coins, &CoingeckoListCoinItem{
			ID:     coin.ID,
			Name:   coin.Name,
			Symbol: coin.Symbol,
		})
	}
	return coins, nil
}

func (cgc *CoinGeckoClient) Ping(ctx context.Context) bool {
	_, err := cgc.c.Ping()
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
