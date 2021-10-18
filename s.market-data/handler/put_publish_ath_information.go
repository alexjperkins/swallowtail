package handler

import (
	"context"
	"fmt"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.market-data/assets"
	"time"

	marketdataproto "swallowtail/s.market-data/proto"

	"github.com/monzo/slog"
)

var (
	athAssets            = assets.LatestPriceAssets
	athTriggerPercentage = 0.025
)

// PublishATHInfo ...
func (s *MarketDataService) PublishATHInfo(
	ctx context.Context, in *marketdataproto.PublishATHInformationRequest,
) (*marketdataproto.PublishATHInformationResponse, error) {
	slog.Trace(ctx, "Market data publishing ATH information")

	for _, asset := range latestPriceAssets {
		asset := asset
		go func() {
			ath, latestPrice, err := fetchATHInfoFromCoingecko(ctx, asset.Symbol, asset.AssetPair)
			if err != nil {
				slog.Error(ctx, "Failed to fetch ATH info from coin gecko: %v: %s %s", err, asset.Symbol, asset.AssetPair)
			}

			if ((ath - latestPrice) / ath) < athTriggerPercentage {
				idempotencyKey := fmt.Sprintf("athinfo-%s-%s-%v-%s", asset.Symbol, asset.AssetPair, ath, time.Now().UTC().Truncate(4*time.Hour))
				if _, err := (&discordproto.SendMsgToChannelRequest{
					Content:        formatATHContent(asset, ath, latestPrice),
					ChannelId:      discordproto.DiscordSatoshiAlertsChannel,
					SenderId:       marketdataproto.MarketDataSystemActor,
					IdempotencyKey: idempotencyKey,
				}).Send(ctx).Response(); err != nil {
					slog.Warn(ctx, "Failed to notifiy discord of approaching ATH for %s%s", asset.Symbol, asset.AssetPair, map[string]string{
						"idempotency_key": idempotencyKey,
					})
				}
			}
		}()
	}

	return &marketdataproto.PublishATHInformationResponse{}, nil
}

func formatATHContent(asset *assets.AssetPair, ath, latestPrice float64) string {
	header := ":robot:     `Approaching ATH: %s%s`     :first_quarter_moon_with_face:"
	content := `

ASSET:        %s%s
ATH:          %v
LATEST PRICE: %v
DIFF:         %v
TIMESTAMP     %v
`
	formattedContent := fmt.Sprintf(content, asset.Symbol, asset.AssetPair, ath, latestPrice, time.Now().UTC().Truncate(15*time.Minute))
	return fmt.Sprintf("%s```%s```", header, formattedContent)
}
