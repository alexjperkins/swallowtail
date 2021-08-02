package handler

import (
	"context"
	"swallowtail/s.coingecko/client"
	coingeckoproto "swallowtail/s.coingecko/proto"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

// GetAssetLatestPriceBySymbol ...
func (s *CoingeckoService) GetAssetLatestPriceBySymbol(
	ctx context.Context, in *coingeckoproto.GetAssetLatestPriceBySymbolRequest,
) (*coingeckoproto.GetAssetLatestPriceBySymbolResponse, error) {
	switch {
	case in.AssetSymbol == "":
		return nil, terrors.PreconditionFailed("missing-param.asset-symbol", "Missing parameter asset symbol", nil)
	case in.AssetPair == "":
		return nil, terrors.PreconditionFailed("missing-param.asset-pair", "Missing parameter asset pair", nil)
	}

	errParams := map[string]string{
		"asset_pair":   in.AssetPair,
		"asset_symbol": in.AssetSymbol,
	}

	slog.Info(ctx, "HERE")

	allTimeHighPrice, err := client.GetCurrentPriceFromSymbol(ctx, in.GetAssetSymbol(), in.GetAssetPair())
	if err != nil {
		return nil, terrors.Augment(err, "Failed to get current price by symbol via coingecko", errParams)
	}

	return &coingeckoproto.GetAssetLatestPriceBySymbolResponse{
		LatestPrice: float32(allTimeHighPrice),
		AssetSymbol: in.AssetSymbol,
	}, nil
}
