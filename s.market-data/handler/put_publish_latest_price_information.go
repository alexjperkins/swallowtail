package handler

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/util"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.market-data/assets"
	marketdataproto "swallowtail/s.market-data/proto"
)

var (
	latestPriceAssets = assets.LatestPriceAssets
)

// AssetInfo ...
type AssetInfo struct {
	Symbol                   string
	AssetPair                string
	LatestPrice              float64
	PriceChangePercentage24h float64
	Group                    string
	ATH                      float64
}

//AssetInfoList ...
type AssetInfoList []*AssetInfo

func (a AssetInfoList) Len() int      { return len(a) }
func (a AssetInfoList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a AssetInfoList) Less(i, j int) bool {
	if a[i].Group == a[j].Group {
		if a[i].Symbol == a[j].Symbol {
			return a[i].AssetPair < a[j].AssetPair
		}

		return a[i].Symbol < a[j].Symbol
	}

	return a[i].Group < a[j].Group
}

// PublishLatestPriceInformation ...
func (s *MarketDataService) PublishLatestPriceInformation(
	ctx context.Context, in *marketdataproto.PublishLatestPriceInformationRequest,
) (*marketdataproto.PublishLatestPriceInformationResponse, error) {
	slog.Trace(ctx, "Market data publishing latest prices")

	var (
		assetInfo = make([]*AssetInfo, 0, len(latestPriceAssets))
		wg        sync.WaitGroup
		mu        sync.RWMutex
	)
	for _, asset := range latestPriceAssets {
		asset := asset

		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(jitter(0, 59))

			latestPrice, change24h, err := fetchLatestPriceFromCoingecko(ctx, asset.Symbol, asset.AssetPair)
			if err != nil {
				slog.Warn(ctx, "Failed to fetch latest price from coingecko: %v: %s%s", err, asset.Symbol, asset.AssetPair)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			assetInfo = append(assetInfo, &AssetInfo{
				Symbol:                   asset.Symbol,
				AssetPair:                asset.AssetPair,
				LatestPrice:              latestPrice,
				PriceChangePercentage24h: change24h,
				Group:                    asset.Grouping.String(),
			})
		}()
	}

	wg.Wait()

	// Sort our asset info alphabetically.
	sort.Slice(assetInfo, func(i, j int) bool {
		switch {
		case assetInfo[i].Group < assetInfo[j].Group:
			return true
		case assetInfo[i].Symbol < assetInfo[j].Symbol:
			return true
		case assetInfo[i].Symbol == assetInfo[j].Symbol:
			return assetInfo[i].AssetPair == assetInfo[j].AssetPair
		default:
			return false
		}
	})

	// TODO: this is the most efficient implementation & we duplicate work here; we should look to improve.
	// But since we poll this endpoint once an hour - it's not the end of the world.
	var (
		indent      int
		priceIndent int
	)
	for _, asset := range assetInfo {
		l := len(fmt.Sprintf("%s%s", asset.Symbol, asset.AssetPair))
		if l > indent {
			indent = l
		}

		pl := len(fmt.Sprintf("%.3f", asset.LatestPrice))
		if pl > priceIndent {
			priceIndent = pl
		}
	}

	// Format content.
	var (
		sb  strings.Builder
		now = time.Now().UTC().Truncate(time.Hour)
	)
	sb.WriteString(fmt.Sprintf(":robot:    `Market Data: Hourly Update: %v`    :dove:\n", now))
	for _, asset := range assetInfo {
		var emoji = ":black_square_large:"
		switch {
		case asset.PriceChangePercentage24h > 0:
			emoji = ":green_square:"
		case asset.PriceChangePercentage24h < 0:
			emoji = ":red_square:"
		}

		sb.WriteString(
			fmt.Sprintf(
				"\n%s `[%s%s]:%s %.3f %s 24h: %.2f%%`",
				emoji,
				strings.ToUpper(asset.Symbol),
				strings.ToUpper(asset.AssetPair),
				strings.Repeat(" ", indent+1+(indent-len(fmt.Sprintf("%s%s", asset.Symbol, asset.AssetPair)))),
				asset.LatestPrice,
				strings.Repeat(" ", priceIndent+(priceIndent-len(fmt.Sprintf("%.3f", asset.LatestPrice)))),
				asset.PriceChangePercentage24h,
			),
		)
	}

	content := sb.String()
	idempotencyKey := fmt.Sprintf("%s-%s-%s", "marketdataprice", util.Sha256Hash(content), now)

	// Publish latest price information to discord.
	if err := batchPublishToDiscord(ctx, sb.String(), discordproto.DiscordSatoshiPriceBotChannel, idempotencyKey); err != nil {
		slog.Error(ctx, "Failed to publish latest price info to discord", map[string]string{
			"idempotency_key": idempotencyKey,
			"error":           err.Error(),
			"len":             strconv.Itoa(len(content)),
		})
	}

	slog.Trace(ctx, "Market data published latest prices")

	return &marketdataproto.PublishLatestPriceInformationResponse{}, nil
}
