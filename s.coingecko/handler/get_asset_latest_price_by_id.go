package handler

import (
	"context"
	"swallowtail/s.coingecko/client"
	coingeckoproto "swallowtail/s.coingecko/proto"

	"github.com/monzo/terrors"
)

// GetAssetLatestPriceByID ...
func (s *CoingeckoService) GetAssetLatestPriceByID(
	ctx context.Context, in *coingeckoproto.GetAssetLatestPriceByIDRequest,
) (*coingeckoproto.GetAssetLatestPriceByIDResponse, error) {
	switch {
	case in.CoingeckoCoinId == "":
		return nil, terrors.PreconditionFailed("missing-param.coingecko-coin-id", "Missing parameter coingecko coin id", nil)
	case in.AssetPair == "":
		return nil, terrors.PreconditionFailed("missing-param.asset-pair", "Missing parameter asset pair", nil)
	}

	errParams := map[string]string{
		"asset_pair":        in.AssetPair,
		"coingecko_coin_id": in.CoingeckoCoinId,
	}

	latestPrice, percentagePriceChange24h, err := client.GetCurrentPriceFromID(ctx, in.GetCoingeckoCoinId(), in.AssetPair)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to get latest price by asset symbol via coingecko", errParams)
	}

	return &coingeckoproto.GetAssetLatestPriceByIDResponse{
		LatestPrice:               float32(latestPrice),
		PercentagePriceChange_24H: float32(percentagePriceChange24h),
		CoingeckoCoinId:           in.CoingeckoCoinId,
	}, nil
}
