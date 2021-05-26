package sync

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"swallowtail/s.googlesheets/domain"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

func formatPagerMsg(name, ticker, errorMsg string, rowIndex int) string {
	return fmt.Sprintf(
		":wave: Hello there, %s\n```Failed to get price for row: `%v` with ticker `%s`\nError: %v\n```Please check it is correct.\n",
		name, strconv.Itoa(rowIndex), ticker, errorMsg,
	)
}

func calculateTotalPNL(rows []*domain.PortfolioRow) float64 {
	var total float64
	for _, row := range rows {
		if row == nil {
			continue
		}
		total += row.PNL
	}
	return total
}

func calculateTotalWorth(ctx context.Context, rows []*domain.PortfolioRow, assetPair string, exchangeClient ExchangeClient) (float64, error) {
	var total float64

	validAssetPair, ok := isValidAssetPairOrConvert(assetPair)
	if !ok {
		return 0.0, terrors.BadRequest("invalid-asset-pair", "Failed to convert asset pair; invalid", map[string]string{
			"asset_pair": "assetPair",
		})
	}

	for _, row := range rows {
		validRowAssetPair, ok := isValidAssetPairOrConvert(row.AssetPair)
		if !ok {
			return 0.0, terrors.BadRequest("invalid-asset-pair", "Failed to convert asset pair; invalid", map[string]string{
				"asset_pair": "assetPair",
			})
		}

		if validRowAssetPair == validAssetPair {
			total += row.CurrentValue
			continue
		}

		slog.Trace(ctx, "Fetching conversion coeff. for %s [%s -> %s]", row.Ticker, row.AssetPair, assetPair)
		coefficient, err := exchangeClient.GetPrice(ctx, validRowAssetPair, validAssetPair)
		if err != nil {
			return 0.0, terrors.Augment(err, "Failed to calculate total net worth", map[string]string{
				"ticker":            row.Ticker,
				"ticker_asset_pair": row.AssetPair,
				"net_asset_pair":    assetPair,
			})
		}
		total += row.CurrentValue * coefficient
	}
	return total, nil
}

func pagerOnIncrease(ctx context.Context, pager func(ctx context.Context, msg string) error, ownerName string, rows []*domain.PortfolioRow, multipleIncreases []float64, delta float64) {
	sort.Float64s(multipleIncreases)
	for _, row := range rows {
		for _, increaseAmount := range multipleIncreases {
			if within(row.AverageEntry, delta) {
				msg := fmt.Sprintf(":wave: Hi %s, %s entry [%v] is close to a multiple target [%v].\nPlease consider taking out your initial investmentl.", ownerName, row.Ticker, row.AverageEntry, increaseAmount)
				if err := pager(ctx, msg); err != nil {
					slog.Error(ctx, "Failed to page %v on increase", ownerName)
				}
			}
		}
	}
}

// calculateHistoricalPNL iterates through all rows, refreshes them to recalculate the PNL
// and returns the total.
func calculateHistoricalPNL(rows []*domain.HistoricalTradeRow) float64 {
	var total float64
	for _, row := range rows {
		row.Refresh()
		total += row.PNL
	}
	return total
}

func isValidAssetPairOrConvert(assetPair string) (string, bool) {
	validAssetPairMtx.RLock()
	defer validAssetPairMtx.RUnlock()
	_, ok := validAssetPairs[assetPair]
	if !ok {
		return "", false
	}
	if strings.ToLower(assetPair) == "usdt" {
		return "usd", true
	}
	return assetPair, true
}

func within(value, delta float64) bool {
	deltaValue := math.Abs(value - (value * delta))
	l, r := value-deltaValue, value+deltaValue
	if value < l {
		return false
	}
	if value > r {
		return false
	}
	return true
}
