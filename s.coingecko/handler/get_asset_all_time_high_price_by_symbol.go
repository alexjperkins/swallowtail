package handler

import (
	"context"

	"github.com/monzo/terrors"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.coingecko/client"
	coingeckoproto "swallowtail/s.coingecko/proto"
)

// GetATHBySymbol ...
func (s *CoingeckoService) GetATHBySymbol(
	ctx context.Context, in *coingeckoproto.GetATHBySymbolRequest,
) (*coingeckoproto.GetATHBySymbolResponse, error) {
	var assetPair = "usd"
	switch {
	case in.AssetSymbol == "":
		return nil, terrors.PreconditionFailed("missing-param.asset-symbol", "Missing parameter asset symbol", nil)
	case in.AssetPair != "":
		assetPair = in.AssetPair
	}

	errParams := map[string]string{
		"asset_pair":   assetPair,
		"asset_symbol": in.AssetSymbol,
	}

	allTimeHighPrice, currentPrice, err := client.GetATHFromSymbol(ctx, in.GetAssetSymbol(), in.AssetPair)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_latest_price_from_coingecko", errParams)
	}

	return &coingeckoproto.GetATHBySymbolResponse{
		AllTimeHighPrice: float32(allTimeHighPrice),
		AssetSymbol:      in.AssetSymbol,
		CurrentPrice:     float32(currentPrice),
	}, nil
}
