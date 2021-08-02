package handler

import (
	"context"
	coingeckoproto "swallowtail/s.coingecko/proto"
	googlesheetsproto "swallowtail/s.googlesheets/proto"

	"github.com/monzo/terrors"
)

func (s *GooglesheetsService) TmpGetLatestPriceBySymbol(
	ctx context.Context, in *googlesheetsproto.TmpGetLatestPriceBySymbolRequest,
) (*googlesheetsproto.TmpGetLatestPriceBySymbolResponse, error) {
	switch {
	case in.AssetPair == "":
		return nil, terrors.PreconditionFailed("missing-param.asset_pair", "Missing param: asset pair", nil)
	case in.AssetSymbol == "":
		return nil, terrors.PreconditionFailed("missing-param.asset_symbol", "Missing param: asset symbol", nil)
	}

	errParams := map[string]string{
		"asset_symbol": in.AssetSymbol,
		"asset_pair":   in.AssetPair,
	}

	rsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
		AssetPair:   in.AssetPair,
		AssetSymbol: in.AssetSymbol,
	}).Send(ctx).Response()
	if err != nil {
		return nil, terrors.Augment(err, "Failed to get the latest price by symbool", errParams)
	}

	return &googlesheetsproto.TmpGetLatestPriceBySymbolResponse{
		LatestPrice: rsp.LatestPrice,
	}, nil
}
