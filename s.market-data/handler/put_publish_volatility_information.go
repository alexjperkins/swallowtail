package handler

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/monzo/slog"

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
		volatilityCache.SetTTL(5 * time.Minute) // Set cache TTL to 5 minute. This is the same default TTL as the price client.
	})

	slog.Trace(ctx, "Market data publishing volatility data.")

	// We have to wait until we've collected all our information for each coin; otherwise gRPC
	// inadvertantly cancels our context.
	var wg sync.WaitGroup
	for _, asset := range volatilityAssets {
		asset := asset
		wg.Add(1)

		go func() {
			defer wg.Done()
			time.Sleep(jitter(0, 59))

			// Fetch the latest price.
			latestPrice, _, err := fetchLatestPriceFromCoingecko(ctx, asset.Symbol, asset.AssetPair)
			if err != nil {
				slog.Error(ctx, "Failed to get latest price for %s%s to determine volatility", asset.Symbol, asset.AssetPair)
				return
			}

			key := fmt.Sprintf("%s%s", asset.Symbol, asset.AssetPair)

			// Fetch previous price from in memory cache.
			previousPrice, err := volatilityCache.Get(key)
			switch {
			case err == ttlcache.ErrNotFound:
				volatilityCache.Set(key, latestPrice)
				return
			case err != nil:
				slog.Error(ctx, "Failed to check for previous price in volatility cache")
				return
			}

			f, ok := previousPrice.(float64)
			if !ok {
				slog.Error(ctx, "Invalid type in volatility cache: expected: float64, got: %T", previousPrice)
				return
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
					"error":           err.Error(),
				})
			}

			if err := volatilityCache.Set(key, latestPrice); err != nil {
				slog.Error(ctx, "Failed to set volatility cache: %s%s", asset.Symbol, asset.AssetPair)
				return
			}
		}()
	}

	// We wait for our goroutines to finish; and also give a 5 second buffer to allow for all
	// discord messages to be posted over the network other the context will be closed.
	wg.Wait()
	<-time.After(5 * time.Second)

	return &marketdataproto.PublishVolatilityInformationResponse{}, nil
}

func formatVolatilityContent(asset *assets.AssetPair, latestPrice, diff float64, increasing bool) string {
	var emoji = ":chart_with_upwards_trend:"
	if !increasing {
		emoji = ":chart_with_downwards_trend:"
	}

	header := fmt.Sprintf(":rotating_light:    `High Volatility Alert: %s%s` %s    :robot:", strings.ToUpper(asset.Symbol), strings.ToUpper(asset.AssetPair), emoji)
	content := `
ASSET:        %s%s
LATEST PRICE: %.3f
15M_CHANGE :  %.2f%%
`
	formattedContent := fmt.Sprintf(content, strings.ToUpper(asset.Symbol), strings.ToUpper(asset.AssetPair), latestPrice, diff*100)
	return fmt.Sprintf("%s```%s```", header, formattedContent)
}
