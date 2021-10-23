package handler

import (
	"context"
	"time"

	"swallowtail/libraries/gerrors"
	binanceproto "swallowtail/s.binance/proto"
	coingeckoproto "swallowtail/s.coingecko/proto"
	discordproto "swallowtail/s.discord/proto"
	ftxproto "swallowtail/s.ftx/proto"
	marketdataproto "swallowtail/s.market-data/proto"

	"github.com/monzo/slog"
)

// publishToDiscord ...
func publishToDiscord(ctx context.Context, content, channel, idempotencyKey string) error {
	if _, err := (&discordproto.SendMsgToChannelRequest{
		Content:        content,
		ChannelId:      channel,
		IdempotencyKey: idempotencyKey,
		SenderId:       marketdataproto.MarketDataSystemActor,
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_publish_msg_to_discord", nil)
	}

	return nil
}

// fetchLatestPriceFromCoingecko ...
func fetchLatestPriceFromCoingecko(ctx context.Context, symbol, assetPair string) (float64, float64, error) {
	rsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
		AssetSymbol: symbol,
		AssetPair:   assetPair,
	}).SendWithTimeout(ctx, 2*time.Minute).Response()
	if err != nil {
		return 0, 0, gerrors.Augment(err, "failed_to_fetch_latest_price_from_coingecko", map[string]string{
			"symbol":     symbol,
			"asset_pair": assetPair,
		})
	}

	return float64(rsp.LatestPrice), float64(rsp.PercentagePriceChange_24H), nil
}

// fetchATHInfoFromCoingecko ...
func fetchATHInfoFromCoingecko(ctx context.Context, symbol, assetPair string) (float64, float64, error) {
	rsp, err := (&coingeckoproto.GetATHBySymbolRequest{
		AssetSymbol: symbol,
		AssetPair:   assetPair,
	}).SendWithTimeout(ctx, 2*time.Minute).Response()
	if err != nil {
		return 0, 0, gerrors.Augment(err, "failed_to_ath_info", map[string]string{
			"symbol":     symbol,
			"asset_pair": assetPair,
		})
	}

	return float64(rsp.AllTimeHighPrice), float64(rsp.CurrentPrice), nil
}

func getFundingRateFromBinance(ctx context.Context, symbol string) (float64, error) {
	rsp, err := (&binanceproto.GetFundingRatesRequest{
		Symbol: symbol,
		Limit:  1,
	}).Send(ctx).Response()
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_get_funding_rate_from_binance", map[string]string{
			"symbol": symbol,
		})
	}

	if len(rsp.FundingRates) == 0 {
		slog.Warn(ctx, "No data for funding rates passed: %s", symbol)
		return 0.0, nil
	}

	return float64(rsp.FundingRates[0].FundingRate), nil
}

func getFundingRateFromFTX(ctx context.Context, symbol string) (float64, error) {
	rsp, err := (&ftxproto.GetFTXFundingRatesRequest{
		Symbol: symbol,
		Limit:  1,
	}).Send(ctx).Response()
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_get_funding_rate_from_ftx", nil)
	}

	if len(rsp.FundingRates) == 0 {
		slog.Warn(ctx, "No data for funding rates passed: %s", symbol)
		return 0.0, nil
	}

	return float64(rsp.FundingRates[0].FundingRate), nil
}
