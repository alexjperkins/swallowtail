package client

import (
	"context"
	"swallowtail/libraries/util"
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

	// Ping serves as a healthcheck to the Binance API.
	Ping(context.Context) error
}

func Init() {
	apiKey := util.SetEnv("BINANCE_API_KEY")
	ctx := context.Background()

	mu.Lock()
	defer mu.Unlock()

	if client != nil {
		return
	}
	c, err := NewDefaultClient(ctx, apiKey)
	if err != nil {
		// Panic since if we can't connect to Binance then this service is as good as dead.
		panic(err)
	}
	client = c
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
