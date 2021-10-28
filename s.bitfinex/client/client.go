package client

import (
	"context"
	"time"

	"github.com/monzo/slog"
	"github.com/opentracing/opentracing-go"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
	"swallowtail/s.bitfinex/dto"
)

var (
	client BitfinexClient
)

// BitfinexClient defines the interface for the Bitfinex Exchange.
type BitfinexClient interface {
	Ping(ctx context.Context) error
	GetStatus(ctx context.Context, req *dto.GetStatusRequest) (*dto.GetStatusResponse, error)
	GetFundingRates(ctx context.Context, req *dto.GetFundingRatesRequest) (*dto.GetFundingRatesResponse, error)
}

// Init initializes the default bitfinex client.
func Init(ctx context.Context) error {
	cli := &bitfinexClient{
		http: transport.NewHTTPClient(10*time.Second, &bitfinexRateLimiter{}),
	}

	if err := cli.Ping(ctx); err != nil {
		panic(gerrors.Augment(err, "failed_to_establish_connection_to_bitfinex", nil))
	}

	slog.Info(ctx, "Established a connection to Bitfinex", nil)

	client = cli
	return nil
}

// GetStatus fetches the status of the bitfinex platfrom.
func GetStatus(ctx context.Context, req *dto.GetStatusRequest) (*dto.GetStatusResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get bitfinex exchange status")
	defer span.Finish()

	then := time.Now().UTC()

	rsp, err := client.GetStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	// Inject server latency into response.
	rsp.ServerLatency = int(time.Since(then) / time.Millisecond)

	return rsp, nil
}

// GetFundingRates ...
func GetFundingRates(ctx context.Context, req *dto.GetFundingRatesRequest) (*dto.GetFundingRatesResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get bitfinex exchange status")
	defer span.Finish()
	return client.GetFundingRates(ctx, req)
}
