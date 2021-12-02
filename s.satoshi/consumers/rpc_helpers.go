package consumers

import (
	"context"
	"swallowtail/libraries/gerrors"
	coingeckoproto "swallowtail/s.coingecko/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
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

func createTradeStrategy(ctx context.Context, tradeStrategy *tradeengineproto.TradeStrategy) (*tradeengineproto.CreateTradeStrategyResponse, error) {
	rsp, err := (&tradeengineproto.CreateTradeStrategyRequest{
		TradeStrategy: tradeStrategy,
	}).Send(ctx).Response()
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
