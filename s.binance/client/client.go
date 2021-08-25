package client

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
	"swallowtail/s.binance/domain"
	"time"

	"github.com/opentracing/opentracing-go"
)

const (
	// Base URL(s)
	binanceAPIUrl  = "https://api.binance.com/api/v3"
	binanceAPIUrl1 = "https://api1.binance.com/api/v3"
	binanceAPIUrl2 = "https://api2.binance.com/api/v3"
	binanceAPIUrl3 = "https://api3.binance.com/api/v3"

	// Base SPOT URL(s)
	binanceSpotURL  = "https://api.binance.com/sapi/v1"
	binanceSpotURL1 = "https://api1.binance.com/sapi/v1"
	binanceSpotURL2 = "https://api2.binance.com/sapi/v1"
	binanceSpotURL3 = "https://api3.binance.com/sapi/v1"
)

var (
	client BinanceClient
)

type BinanceClient interface {
	// ListAllAssetPairs makes a call to Binance to retrieve all the futures tradable asset pairs.
	ListAllAssetPairs(context.Context) (*ListAllAssetPairsResponse, error)

	// ExecuteSpotTrade attempts to execute a spot trade on Binance.
	ExecuteSpotTrade(ctx context.Context, trade *domain.Trade) error

	// Ping serves as a healthcheck to the Binance API.
	Ping(context.Context) error

	// ReadSpotAccount reads from the users spot account.
	ReadSpotAccount(context.Context, *ReadSpotAccountRequest) (*ReadSpotAccountResponse, error)

	// ReadPerpetualFuturesAccount reads from the users perpetual futures account.
	ReadPerpetualFuturesAccount(context.Context, *ReadPerptualFuturesAccountRequest) (*ReadPerptualFuturesAccountResponse, error)

	// VerifyCredentials verifies the given credentials of the users.
	VerifyCredentials(context.Context, *Credentials) (*VerifyCredentialsResponse, error)
}

// Init initializes the default binance client for this service.
func Init(ctx context.Context) error {
	c := &binanceClient{
		http: transport.NewHTTPClient(30 * time.Second),
	}

	if err := c.Ping(ctx); err != nil {
		// Panic since if we can't connect to Binance then this service is as good as dead.
		return gerrors.Augment(err, "failed.binance_client_initialization", nil)
	}

	client = c
	return nil
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

// ReadSpotAccount ...
func ReadSpotAccount(ctx context.Context, req *ReadSpotAccountRequest) (*ReadSpotAccountResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Read from binance spot account")
	defer span.Finish()
	return nil, nil
}

// ReadPerpetualFuturesAccount ...
func ReadPerpetualFuturesAccount(ctx context.Context, req *ReadPerptualFuturesAccountRequest) (*ReadPerptualFuturesAccountResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Read from binance perpetual futures account")
	defer span.Finish()
	return nil, nil
}

// VerifyCredentials ...
func VerifyCredentials(ctx context.Context, credentials *Credentials) (*VerifyCredentialsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Verify credentials for user")
	defer span.Finish()
	return client.VerifyCredentials(ctx, credentials)
}
