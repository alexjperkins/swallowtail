package handler

import (
	"context"
	"swallowtail/s.coingecko/client"
	coingeckoproto "swallowtail/s.coingecko/proto"

	"github.com/monzo/terrors"
)

// GetATHBySymbol ...
func (s *CoingeckoService) GetATHBySymbol(
	ctx context.Context, in *coingeckoproto.GetATHBySymbolRequest,
) (*coingeckoproto.GetATHBySymbolResponse, error) {
	switch {
	case in.AssetSymbol == "":
		return nil, terrors.PreconditionFailed("missing-param.asset-symbol", "Missing parameter asset symbol", nil)
	}

	errParams := map[string]string{
		"asset_pair":   in.AssetPair,
		"asset_symbol": in.AssetSymbol,
	}

	allTimeHighPrice, err := client.GetATHFromSymbol(ctx, in.GetAssetSymbol())
	if err != nil {
		return nil, terrors.Augment(err, "Failed to get all time high price by symbol via coingecko", errParams)
	}

	return &coingeckoproto.GetATHBySymbolResponse{
		AllTimeHighPrice: float32(allTimeHighPrice),
		AssetSymbol:      in.AssetSymbol,
	}, nil
}
