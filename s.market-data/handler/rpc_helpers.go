package handler

import (
	"context"
	"time"

	"swallowtail/libraries/gerrors"
	coingeckoproto "swallowtail/s.coingecko/proto"
	discordproto "swallowtail/s.discord/proto"
	marketdataproto "swallowtail/s.market-data/proto"
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
func fetchLatestPriceFromCoingecko(ctx context.Context, symbol, assetPair string) (float64, error) {
	rsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
		AssetSymbol: symbol,
		AssetPair:   assetPair,
	}).SendWithTimeout(ctx, 2*time.Minute).Response()
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_fetch_latest_price_from_coingecko", map[string]string{
			"symbol":     symbol,
			"asset_pair": assetPair,
		})
	}

	return float64(rsp.LatestPrice), nil
}
