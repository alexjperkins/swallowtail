package bybt

import "context"

type bybtClient struct{}

func (b *bybtClient) GetExchangeFundingRatesByAsset(
	ctx context.Context, in *GetExchangeFundingRatesByAssetRequest,
) (*GetExchangeFundingRatesByAssetResponse, error) {
	return nil, nil
}

func (b *bybtClient) Ping(ctx context.Context) error {
	return nil
}
