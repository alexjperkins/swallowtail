package consumers

import (
	"context"
	"swallowtail/libraries/gerrors"
	coingeckoproto "swallowtail/s.coingecko/proto"
)

func getAssetLatestPrice(ctx context.Context, symbol, assetPair string) (float64, error) {
	rsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
		AssetSymbol: symbol,
		AssetPair:   assetPair,
	}).Send(ctx).Response()
	if err != nil {
		return 0, gerrors.Augment(err, "failed_get_asset_latest_price", map[string]string{
			"asset_symbol": symbol,
			"assetPair":    assetPair,
		})
	}
	return float64(rsp.LatestPrice), nil
}