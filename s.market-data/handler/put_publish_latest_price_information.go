package handler

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"swallowtail/libraries/util"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.market-data/assets"
	marketdataproto "swallowtail/s.market-data/proto"

	"github.com/monzo/slog"
)

var (
	latestPriceAssets = assets.LatestPriceAssets
)

// LatestPriceInfo ...
type LatestPriceInfo struct {
	Symbol                   string
	AssetPair                string
	LatestPrice              float64
	PriceChangePercentage24h float64
}

// PublishLatestPriceInformation ...
func (s *MarketDataService) PublishLatestPriceInformation(
	ctx context.Context, in *marketdataproto.PublishLatestPriceInformationRequest,
) (*marketdataproto.PublishLatestPriceInformationResponse, error) {
	slog.Trace(ctx, "Market data publishing latest prices")

	var (
		assetInfo = make([]*LatestPriceInfo, 0, len(latestPriceAssets))
		wg        sync.WaitGroup
		mu        sync.RWMutex
	)
	for _, asset := range latestPriceAssets {
		asset := asset

		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(jitter(0, 119))

			latestPrice, err := fetchLatestPriceFromCoingecko(ctx, asset.Symbol, asset.AssetPair)
			if err != nil {
				slog.Warn(ctx, "Failed to fetch latest price from coingecko: %v: %s%s", err, asset.Symbol, asset.AssetPair)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			assetInfo = append(assetInfo, &LatestPriceInfo{
				Symbol:      asset.Symbol,
				AssetPair:   asset.AssetPair,
				LatestPrice: latestPrice,
			})
		}()
	}

	wg.Wait()

	// Sort our asset info alphabetically.
	sort.Slice(assetInfo, func(i, j int) bool {
		switch {
		case assetInfo[i].Symbol < assetInfo[j].Symbol:
			return true
		case assetInfo[i].Symbol == assetInfo[j].Symbol:
			return assetInfo[i].AssetPair == assetInfo[j].AssetPair
		default:
			return false
		}
	})

	var indent int
	for _, asset := range assetInfo {
		l := len(fmt.Sprintf("%s%s", asset.Symbol, asset.AssetPair))
		if l > indent {
			indent = l
		}
	}

	var (
		sb  strings.Builder
		now = time.Now().UTC().Truncate(time.Hour)
	)

	// Format content.
	sb.WriteString(fmt.Sprintf(":robot:    Market Data: Hourly Update: %v    :dove:\n", now))
	for _, asset := range assetInfo {
		var emoji = ":black_square_large:"
		switch {
		case asset.PriceChangePercentage24h > 0:
			emoji = ":green_square:"
		case asset.PriceChangePercentage24h < 0:
			emoji = ":red_square:"
		}

		sb.WriteString(fmt.Sprintf("\n%s%s %s:%s %v 24h: %s%%", asset.Symbol, asset.AssetPair, emoji, strings.Repeat(" ", indent+1), asset.LatestPrice, asset.PriceChangePercentage24h))
	}

	content := sb.String()
	idempotencyKey := fmt.Sprintf("%s-%s-%s", "marketdataprice", util.Sha256Hash(content), now)

	// Publish latest price information to discord.
	if err := publishToDiscord(ctx, sb.String(), discordproto.DiscordSatoshiPriceBotChannel, idempotencyKey); err != nil {
		slog.Error(ctx, "Failed to publish latest price info to discord", map[string]string{
			"idempotency_key": idempotencyKey,
		})
	}

	return &marketdataproto.PublishLatestPriceInformationResponse{}, nil
}
