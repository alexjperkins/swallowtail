package sync

import (
	"context"
	coingeckoproto "swallowtail/s.coingecko/proto"
	"time"

	"github.com/monzo/terrors"
)

func getlatestPriceByID(ctx context.Context, assetID, assetPair string) (*coingeckoproto.GetAssetLatestPriceByIDResponse, error) {
	rsp, err := (&coingeckoproto.GetAssetLatestPriceByIDRequest{
		CoingeckoCoinId: assetID,
		AssetPair:       assetPair,
	}).SendWithTimeout(ctx, 30*time.Second).Response()
	if err != nil {
		return nil, terrors.Augment(err, "Failed to get the latest price", map[string]string{
			"asset_id":   assetID,
			"asset_pair": assetPair,
		})

	}

	return rsp, nil
}

func getLatestPriceBySymbol(ctx context.Context, symbol, assetPair string) (*coingeckoproto.GetAssetLatestPriceBySymbolResponse, error) {
	rsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
		AssetSymbol: symbol,
		AssetPair:   assetPair,
	}).SendWithTimeout(ctx, 30*time.Second).Response()
	if err != nil {
		return nil, terrors.Augment(err, "Failed to get the latest price by symbol", map[string]string{
			"asset_symbol": symbol,
			"asset_pair":   assetPair,
		})
	}

	return rsp, nil
}
