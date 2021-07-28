package bybt

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

var (
	client ByBtClient
)

// ByBtClient ...
type ByBtClient interface {
	GetExchangeFundingRatesByAsset(ctx context.Context, in *GetExchangeFundingRatesByAssetRequest) (*GetExchangeFundingRatesByAssetResponse, error)
	Ping(ctx context.Context) error
}

// Init ...
func Init(ctx context.Context) error {
	client = &bybtClient{}
	return client.Ping(ctx)
}

// GetExchangeFundingRatesByAsset ...
func GetExchangeFundingRatesByAsset(
	ctx context.Context, in *GetExchangeFundingRatesByAssetRequest,
) (*GetExchangeFundingRatesByAssetResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Send bybt request exchange funding rates")
	defer span.Finish()
	return nil, nil
}
