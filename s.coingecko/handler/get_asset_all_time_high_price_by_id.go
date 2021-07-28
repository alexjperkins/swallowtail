package handler

import (
	"context"
	"swallowtail/s.coingecko/client"
	coingeckoproto "swallowtail/s.coingecko/proto"

	"github.com/monzo/terrors"
)

// GetATHByID ...
func (s *CoingeckoService) GetATHByID(
	ctx context.Context, in *coingeckoproto.GetATHByIDRequest,
) (*coingeckoproto.GetATHByIDResponse, error) {
	switch {
	case in.CoingeckoCoinId == "":
		return nil, terrors.PreconditionFailed("missing-param.coingecko-coin-id", "Missing parameter coingecko coin id", nil)
	}

	errParams := map[string]string{
		"asset_pair":        in.AssetPair,
		"coingecko_coin_id": in.CoingeckoCoinId,
	}

	allTimeHighPrice, err := client.GetATHFromID(ctx, in.GetCoingeckoCoinId())
	if err != nil {
		return nil, terrors.Augment(err, "Failed to get all time high price by id via coingecko", errParams)
	}

	return &coingeckoproto.GetATHByIDResponse{
		AllTimeHighPrice: float32(allTimeHighPrice),
		CoingeckoCoinId:  in.CoingeckoCoinId,
	}, nil
}
