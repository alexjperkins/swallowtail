package handler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/monzo/slog"

	coingeckoproto "swallowtail/s.coingecko/proto"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.market-data/assets"
	marketdataproto "swallowtail/s.market-data/proto"
)

var (
	volatilityAssets = assets.LatestPriceAssets
	volatilityCache  *ttlcache.Cache
	volatilityOnce   sync.Once
)

// PublishVolatilityInformation ...
func (s *MarketDataService) PublishVolatilityInformation(
	ctx context.Context, in *marketdataproto.PublishVolatilityInformationRequest,
) (*marketdataproto.PublishVolatilityInformationResponse, error) {
	volatilityOnce.Do(func() {
		volatilityCache = ttlcache.NewCache()
		volatilityCache.SetCacheSizeLimit(len(volatilityAssets))
		volatilityCache.SetTTL(20 * time.Minute) // Whilst we're running on a 15 minute cron, we add 5 mins to act as buffer.
	})

	for _, asset := range volatilityAssets {
		asset := asset

		go func() {
			rsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
				AssetSymbol: asset.Symbol,
				AssetPair:   asset.AssetPair,
			}).SendWithTimeout(ctx, 2*time.Minute).Response()
			if err != nil {
				slog.Error(ctx, "Failed to get latest price for %s%s to determine volatility", asset.Symbol, asset.AssetPair)
			}

			key := fmt.Sprintf("%s%s", asset.Symbol, asset.AssetPair)
			latestPrice := float64(rsp.LatestPrice)

			previousPrice, err := volatilityCache.Get(key)
			switch {
			case err == ttlcache.ErrNotFound:
				volatilityCache.Set(key, latestPrice)
			case err != nil:
				slog.Error(ctx, "Failed to check for previous price in volatility cache")
				return
			}

			f, ok := previousPrice.(float64)
			if !ok {
				slog.Error(ctx, "Invalid type in volatility cache: expected: float64, got: %T", previousPrice)
			}

			diff := (latestPrice - f) / f

			var increasing bool
			switch {
			case abs(diff) < asset.VolatilityRating.PercentageTriggerValue():
				// No volatility so skipping.
				return
			case diff > 0:
				increasing = true
			case diff < 0:
				increasing = false
			}

			idempotencyKey := fmt.Sprintf("volinfo-%s-%s-%s", asset.Symbol, asset.AssetPair, time.Now().UTC().Truncate(10*time.Minute))
			if _, err := (&discordproto.SendMsgToChannelRequest{
				Content:        formatVolatilityContent(asset, latestPrice, diff, increasing),
				IdempotencyKey: idempotencyKey,
				ChannelId:      discordproto.DiscordSatoshiAlertsChannel,
				SenderId:       marketdataproto.MarketDataSystemActor,
			}).Send(ctx).Response(); err != nil {
				slog.Error(ctx, "Failed to notify discord of volatility info for: %s%s", asset.Symbol, asset.AssetPair, map[string]string{
					"idempotency_key": idempotencyKey,
				})
			}
		}()
	}

	return &marketdataproto.PublishVolatilityInformationResponse{}, nil
}

func formatVolatilityContent(asset *assets.AssetPair, latestPrice, diff float64, increasing bool) string {
	var emoji = ":chart_with_upwards_trend"
	if !increasing {
		emoji = ":chart_with_downwards_trend:"
	}

	header := fmt.Sprintf(":rotating_light:    `High Volatility Alert: %s: %s%s`    :robot:", emoji, asset.Symbol, asset.AssetPair)
	content := `
ASSET:        %s%s
LATEST PRICE: %v
15M_CHANGE :  %v%%
`
	formattedContent := fmt.Sprintf(content, asset.Symbol, asset.AssetPair, latestPrice, diff*100)
	return fmt.Sprintf("%s```%s```", header, formattedContent)
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}

	return a
}
