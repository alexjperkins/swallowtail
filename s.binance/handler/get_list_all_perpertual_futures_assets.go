package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	binanceproto "swallowtail/s.binance/proto"
)

// ListAllAssetPairs ...
func (s *BinanceService) ListAllAssetPairs(
	ctx context.Context, in *binanceproto.ListAllAssetPairsRequest,
) (*binanceproto.ListAllAssetPairsResponse, error) {
	assetPairs, err := client.ListAllAssetPairs(ctx)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_all_asset_pairs", nil)
	}

	protos := []*binanceproto.AssetPair{}
	for _, ap := range assetPairs.Symbols {
		protos = append(protos, &binanceproto.AssetPair{
			Symbol:            ap.Symbol,
			BaseAsset:         ap.BaseAsset,
			WithMarginTrading: ap.WithMarginTrading,
			WithSpotTrading:   ap.WithSpotTrading,
		})
	}

	return &binanceproto.ListAllAssetPairsResponse{
		AssetPairs: protos,
	}, nil
}
