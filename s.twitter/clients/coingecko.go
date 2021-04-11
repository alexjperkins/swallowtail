package clients

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/monzo/slog"
	coingecko "github.com/superoo7/go-gecko/v3"
)

var (
	CoingeckoClientID = "coingecko-client-id"

	defaultCoinGeckoClientTimeout = time.Second * 10
)

func NewCoinGeckoClient() *CoinGeckoClient {
	httpClient := &http.Client{
		Timeout: defaultCoinGeckoClientTimeout,
	}
	cgc := &CoinGeckoClient{
		c: coingecko.NewClient(httpClient),
	}
	if ok := cgc.Ping(); ok {
		slog.Info(context.TODO(), "Failed to connect coingecko client")
	}
	slog.Info(context.TODO(), "Created coingecko client")
	return cgc
}

type CoinGeckoClient struct {
	c *coingecko.Client
}

func (cgc *CoinGeckoClient) GetATHFromID(id string) (float64, error) {
	coinID, err := cgc.c.CoinsID(strings.ToLower(id), true, false, true, false, false, false)
	if err != nil {
		return 0.0, err
	}

	return coinID.MarketData.ATH["usd"], nil
}

func (cgc *CoinGeckoClient) GetCurrentPriceFromID(id string) (float64, error) {
	ssp, err := cgc.c.SimpleSinglePrice(strings.ToLower(id), "usd")
	if err != nil {
		return 0.0, nil
	}
	return float64(ssp.MarketPrice), nil
}

func (cgc *CoinGeckoClient) Ping() bool {
	_, err := cgc.c.Ping()
	if err != nil {
		return false
	}
	return true
}
