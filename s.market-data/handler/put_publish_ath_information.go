package handler

import (
	"context"
	"fmt"
	"strings"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.market-data/assets"
	"sync"
	"time"

	marketdataproto "swallowtail/s.market-data/proto"

	"github.com/monzo/slog"
)

var (
	athAssets            = assets.LatestPriceAssets
	athTriggerPercentage = 0.025 // 2.5%
)

// PublishATHInformation ...
func (s *MarketDataService) PublishATHInformation(
	ctx context.Context, in *marketdataproto.PublishATHInformationRequest,
) (*marketdataproto.PublishATHInformationResponse, error) {
	slog.Trace(ctx, "Market data publishing ATH information")

	var wg sync.WaitGroup
	for _, asset := range latestPriceAssets {
		asset := asset
		wg.Add(1)

		go func() {
			defer wg.Done()

			ath, latestPrice, err := fetchATHInfoFromCoingecko(ctx, asset.Symbol, asset.AssetPair)
			if err != nil {
				slog.Error(ctx, "Failed to fetch ATH info from coin gecko: %v: %s %s", err, asset.Symbol, asset.AssetPair)
			}

			diff := abs((ath - latestPrice) / ath)
			if diff < athTriggerPercentage {
				// Idempotent on the ATH price, symbol & asset pair.
				idempotencyKey := fmt.Sprintf("athinfo-%s-%s-%vs", asset.Symbol, asset.AssetPair, ath)
				if _, err := (&discordproto.SendMsgToChannelRequest{
					Content:        formatApproachingATHContent(asset, ath, latestPrice, diff),
					ChannelId:      discordproto.DiscordSatoshiAlertsChannel,
					SenderId:       marketdataproto.MarketDataSystemActor,
					IdempotencyKey: idempotencyKey,
				}).Send(ctx).Response(); err != nil {
					slog.Warn(ctx, "Failed to notifiy discord of approaching ATH for %s%s", asset.Symbol, asset.AssetPair, map[string]string{
						"idempotency_key": idempotencyKey,
						"error":           err.Error(),
					})
				}
			}
		}()
	}

	// We wait for our goroutines to finish; and also give a 5 second buffer to allow for all
	// discord messages to be posted over the network other the context will be closed.
	wg.Wait()
	<-time.After(5 * time.Second)

	return &marketdataproto.PublishATHInformationResponse{}, nil
}

func formatApproachingATHContent(asset *assets.AssetPair, ath, latestPrice, diff float64) string {
	header := fmt.Sprintf(":robot:     `Approaching ATH: %s%s`    :first_quarter_moon_with_face:", strings.ToUpper(asset.Symbol), strings.ToUpper(asset.AssetPair))
	content := `

ASSET:        %s%s
ATH:          %.3f
LATEST PRICE: %.3f
DIFF:         %.2f%%
TIMESTAMP     %v
`
	formattedContent := fmt.Sprintf(content, strings.ToUpper(asset.Symbol), strings.ToUpper(asset.AssetPair), ath, latestPrice, diff*100, time.Now().UTC().Truncate(15*time.Minute))
	return fmt.Sprintf("%s```%s```", header, formattedContent)
}
