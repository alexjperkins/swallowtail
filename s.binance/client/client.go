package client

import (
	"context"
	"swallowtail/libraries/util"
	"swallowtail/s.binance/domain"
	"sync"

	"github.com/opentracing/opentracing-go"
)

const (
	binanceAPIUrl = "https://api.binance.com/api/v3"
)

var (
	client BinanceClient
	mu     sync.Mutex
)

type BinanceClient interface {
	// ListAllAssetPairs makes a call to Binance to retrieve all the futures tradable asset pairs.
	ListAllAssetPairs(context.Context) (*ListAllAssetPairsResponse, error)
	// ExecuteSpotTrade attempts to execute a spot trade on Binance.
	ExecuteSpotTrade(ctx context.Context, trade *domain.Trade) error
	// Ping serves as a healthcheck to the Binance API.
	Ping(context.Context) error
}

func Init(ctx context.Context) error {
	apiKey := util.SetEnv("BINANCE_API_KEY")

	mu.Lock()
	defer mu.Unlock()

	if client != nil {
		return nil
	}
	c, err := NewDefaultClient(ctx, apiKey)
	if err != nil {
		// Panic since if we can't connect to Binance then this service is as good as dead.
		return err
	}
	client = c
	return nil
}

func UseMock() {
	mu.Lock()
	defer mu.Unlock()
	client = &mockClient{}
}

// ListAllAssetPairs forwards the response of the binance client; it also adds opentracing span to the
// to the context of the request.
func ListAllAssetPairs(ctx context.Context) (*ListAllAssetPairsResponse, error) {
	// TODO: add timing metrics.
	span, ctx := opentracing.StartSpanFromContext(ctx, "List all Binance asset pairs")
	defer span.Finish()
	return client.ListAllAssetPairs(ctx)
}

// ExecuteSpotTrade ...
func ExecuteSpotTrade(ctx context.Context, trade *domain.Trade) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Execute binance spot trade")
	defer span.Finish()
	return nil
}
