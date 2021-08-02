package sync

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"swallowtail/s.googlesheets/domain"

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

func calculateTotalWorth(ctx context.Context, rows []*domain.PortfolioRow, assetPair string) (float64, error) {
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
			return 0, terrors.BadRequest("invalid-asset-pair", "Failed to convert asset pair; invalid", map[string]string{
				"asset_pair": "assetPair",
			})
		}

		if validRowAssetPair == validAssetPair {
			total += row.CurrentValue
			continue
		}

		rsp, err := getLatestPriceBySymbol(ctx, validRowAssetPair, validAssetPair)
		if err != nil {
			return 0, terrors.Augment(err, "Failed to calculate total worth; error fetching conversion rates", map[string]string{
				"base":    validRowAssetPair,
				"counter": validAssetPair,
			})
		}

		coefficient := float64(rsp.LatestPrice)
		total += row.CurrentValue * coefficient
	}

	return total, nil
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
