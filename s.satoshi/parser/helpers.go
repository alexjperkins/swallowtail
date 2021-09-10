package parser

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dlclark/regexp2"

	"swallowtail/libraries/gerrors"
	coingeckoproto "swallowtail/s.coingecko/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

var (
	// Give me:
	// - All the numbers in a string
	// - That are either integers or decimals of any precision, that are not percentages and are not RR valuess
	numeric = regexp2.MustCompile(`(\b\d+(?:[\.,]\d+)?\b(?!(?:[\.,]\d+)|(?:\s*(?:%|percent|RR|hr|\$|\.|\/))))`, regexp2.None)

	// TODO: we need proper gRPC service testing.
	fetchLatestPrice = getLatestPrice
)

func parseNumbersFromContent(content string) ([]float64, error) {
	var err error
	match, err := numeric.FindStringMatch(content)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_capture_numbers", nil)
	}

	floats := []float64{}
	for {
		if match == nil {
			break
		}

		// Remove any spaces.
		trimmed := strings.ReplaceAll(match.String(), " ", "")
		if trimmed == "" {
			continue
		}

		// Remove any commas e.g 50,000
		trimmed = strings.ReplaceAll(trimmed, ",", "")

		f, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_capture_numbers.failed_to_parse_float", map[string]string{
				"num_str": trimmed,
			})
		}

		floats = append(floats, f)

		match, err = numeric.FindNextMatch(match)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_capture_numbers", nil)
		}
	}

	return floats, nil
}

func parseSide(content string) (tradeengineproto.TRADE_SIDE, bool) {
	fields := strings.Fields(content)

	for _, f := range fields {
		switch strings.ToLower(f) {
		case "long":
			return tradeengineproto.TRADE_SIDE_BUY, true
		case "short":
			return tradeengineproto.TRADE_SIDE_SELL, true
		}
	}

	// We default to longing. It is crypto after all.
	return tradeengineproto.TRADE_SIDE_BUY, false
}

// containsTicker checks if the contain contains a ticker that is traded on Binance
// it assumes that the content passed with be normalized to lowercase.
func parseTicker(content string) string {
	tokens := strings.Fields(strings.ToLower(content))
	for _, token := range tokens {
		switch {
		case
			token == "usd",
			token == "usdt",
			token == "usdc":
			// If we match against some stablecoin inadvertly; then we can skip.
			continue
		case
			strings.Contains(token, "usd"),
			strings.Contains(token, "usdc"),
			strings.Contains(token, "usdt"):
			// But if a token contains a stable coins, then lets assume it's of the form BTCUSDT.
			// We might pick up typos and similar here, but that's fine for now.
			return token

		case strings.Contains(token, "/"):
			// Some mods format their trades as `BTC/USDT`.
			childContent := strings.ReplaceAll(token, "/", " ")
			if parseTicker(childContent) != "" {
				return token
			}
		}

		if _, ok := binanceAssetPairs[token]; ok {
			return token
		}
	}

	return ""
}

func withinRange(value, truth, rangeAsPercentage float64) bool {
	boundary := truth * rangeAsPercentage * 0.01

	switch {
	case value <= truth-boundary:
		return false
	case value >= truth+boundary:
		return false
	default:
		return true
	}
}

func getLatestPrice(ctx context.Context, asset, pair string) (float64, error) {
	rsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
		AssetPair:   pair,
		AssetSymbol: asset,
	}).SendWithTimeout(ctx, 15*time.Second).Response()
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_get_latest_price", nil)
	}

	return float64(rsp.LatestPrice), nil
}

func validatePosition(asset, pair string, truth float64, possibleValues []float64) (entry float64, stopLoss float64, takeProfits []float64, err error) {
	switch {
	case len(possibleValues) < 2:
		return 0, 0, nil, gerrors.FailedPrecondition("failed_to_validate_position.not_enough_values", nil)
	case len(possibleValues) > 1:
		stopLoss = possibleValues[0]
		entry = possibleValues[1]
		fallthrough
	case len(possibleValues) > 2:
		takeProfits = possibleValues[2:]
	}

	errParams := map[string]string{
		"asset":     asset,
		"pair":      pair,
		"stop_loss": fmt.Sprintf("%v", stopLoss),
		"entry":     fmt.Sprintf("%v", entry),
	}

	if !withinRange(stopLoss, truth, 50) {
		// If the first value we parse is way off; then we pop it off and attempt on the rest of the values.
		// The assumption is that for any trade the entry, stop loss & take profits should be relatively close to
		// the current value for perps
		return validatePosition(asset, pair, truth, possibleValues[1:])
	}

	if !withinRange(entry, truth, 50) {
		return 0, 0, nil, gerrors.FailedPrecondition("failed_to_validate_position.bad_stop_loss", errParams)
	}

	validTakeProfits := []float64{}
	for _, tp := range takeProfits {
		if withinRange(tp, truth, 300) {
			validTakeProfits = append(validTakeProfits, tp)
		}
	}

	return entry, stopLoss, validTakeProfits, nil
}

func parseOrderType(content string, currentValue, entry float64, side tradeengineproto.TRADE_SIDE) (tradeengineproto.ORDER_TYPE, bool) {
	var containsLimit bool
	for _, f := range strings.Fields(strings.ToLower(content)) {
		if f == "limit" {
			containsLimit = true
		}
	}

	switch {
	case !containsLimit:
		return tradeengineproto.ORDER_TYPE_MARKET, true
	case side == tradeengineproto.TRADE_SIDE_BUY && entry < currentValue:
		return tradeengineproto.ORDER_TYPE_LIMIT, true
	case side == tradeengineproto.TRADE_SIDE_SELL && entry > currentValue:
		return tradeengineproto.ORDER_TYPE_LIMIT, true
	default:
		return tradeengineproto.ORDER_TYPE_MARKET, false
	}
}
