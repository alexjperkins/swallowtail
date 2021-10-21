package client

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
	"swallowtail/libraries/util"
	"swallowtail/s.binance/domain"
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

	// Base FUTURES URL(s)
	binanceFuturesURL   = "https://fapi.binance.com/fapi/v1"
	binanceFuturesURLV2 = "https://fapi.binance.com/fapi/v2"
)

var (
	client             BinanceClient
	defaultCredentials *Credentials
)

// BinanceClient defines the contract for connection to the Binance Exchange API.
type BinanceClient interface {
	// GetLatestPrices gets all the latest prices from the passed symbols.
	GetLatestPrice(context.Context, *GetLatestPriceRequest) (*GetLatestPriceResponse, error)

	// ListAllAssetPairs makes a call to Binance to retrieve all the futures tradable asset pairs.
	ListAllAssetPairs(context.Context) (*ListAllAssetPairsResponse, error)

	// ExecuteSpotTrade attempts to execute a spot trade on Binance.
	ExecuteSpotTrade(ctx context.Context, trade *domain.Trade) error

	// ExecutePerpetualFuturesTrade
	ExecutePerpetualFuturesTrade(ctx context.Context, req *ExecutePerpetualFuturesTradeRequest, credentials *Credentials) (*ExecutePerpetualFuturesTradeResponse, error)

	// Ping serves as a healthcheck to the Binance API.
	Ping(context.Context) error

	// ReadSpotAccount reads from the users spot account.
	ReadSpotAccount(context.Context, *ReadSpotAccountRequest) (*ReadSpotAccountResponse, error)

	// ReadPerpetualFuturesAccount reads from the users perpetual futures account.
	ReadPerpetualFuturesAccount(context.Context, *ReadPerpetualFuturesAccountRequest, *Credentials) (*ReadPerpetualFuturesAccountResponse, error)

	// GetFuturesExchangeInfo returns the exchange information specifically for perpetual futures contracts.
	GetFuturesExchangeInfo(context.Context, *GetFuturesExchangeInfoRequest) (*GetFuturesExchangeInfoResponse, error)

	GetFundingRate(context.Context, *GetFundingRateRequest) (*GetFundingRateResponse, error)

	// VerifyCredentials verifies the given credentials of the users.
	VerifyCredentials(context.Context, *Credentials) (*VerifyCredentialsResponse, error)

	// GetStatus returns the statistics round the exchange server time & latency.
	GetStatus(context.Context) (*GetStatusResponse, error)
}

// Init initializes the default binance client for this service.
func Init(ctx context.Context) error {
	c := &binanceClient{
		http: transport.NewHTTPClient(30*time.Second, &binanceRateLimiter{}),
	}

	if err := c.Ping(ctx); err != nil {
		// Panic since if we can't connect to Binance then this service is as good as dead.
		return gerrors.Augment(err, "failed.binance_client_initialization", nil)
	}

	apiKey := util.SetEnv("BINANCE_DEFAULT_API_KEY")
	secretKey := util.SetEnv("BINANCE_DEFAULT_SECRET_KEY")

	if apiKey == "" || secretKey == "" {
		return gerrors.FailedPrecondition("failed_to_init_binance_client.credentials_not_set", nil)
	}

	defaultCredentials = &Credentials{
		APIKey:    apiKey,
		SecretKey: secretKey,
	}

	client = c
	return nil
}

// GetLatestPrices ...
func GetLatestPrice(ctx context.Context, req *GetLatestPriceRequest) (*GetLatestPriceResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get all latest prices")
	defer span.Finish()
	return client.GetLatestPrice(ctx, req)
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
func ReadPerpetualFuturesAccount(ctx context.Context, req *ReadPerpetualFuturesAccountRequest, credentials *Credentials) (*ReadPerpetualFuturesAccountResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Read from binance perpetual futures account")
	defer span.Finish()
	return client.ReadPerpetualFuturesAccount(ctx, req, credentials)
}

// ExecutePerpetualFuturesTrade ...
func ExecutePerpetualFuturesTrade(ctx context.Context, req *ExecutePerpetualFuturesTradeRequest, credentials *Credentials) (*ExecutePerpetualFuturesTradeResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Execute perpetual futures trade on binance")
	defer span.Finish()
	return client.ExecutePerpetualFuturesTrade(ctx, req, credentials)
}

// VerifyCredentials ...
func VerifyCredentials(ctx context.Context, credentials *Credentials) (*VerifyCredentialsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Verify credentials for user")
	defer span.Finish()
	return client.VerifyCredentials(ctx, credentials)
}

// GetStatus ...
func GetStatus(ctx context.Context) (*GetStatusResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get Binance exchange status")
	defer span.Finish()

	requestStart := time.Now().UTC()

	rsp, err := client.GetStatus(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	// This is not entirely accurate since it includes the time take for the server to respond
	// This would be ~ systematic - but a more slightly accurate approach would be the subtract half the latency (again not great).
	assumedClockDrift := now.Sub(time.Unix(int64(rsp.ServerTime), 0))

	requestEnd := now.Sub(requestStart)
	rsp.ServerLatency = requestEnd
	rsp.AssumedClockDrift = assumedClockDrift

	return rsp, nil
}

// GetFuturesExchangeInfo ...
func GetFuturesExchangeInfo(ctx context.Context, req *GetFuturesExchangeInfoRequest) (*GetFuturesExchangeInfoResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get Binance futures exchange info")
	defer span.Finish()
	return client.GetFuturesExchangeInfo(ctx, req)
}

// GetFundingRate ...
func GetFundingRate(ctx context.Context, req *GetFundingRateRequest) (*GetFundingRateResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get Binance futures funding rate")
	defer span.Finish()
	return client.GetFundingRate(ctx, req)
}
