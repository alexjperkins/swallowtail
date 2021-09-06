package parser

import (
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"

	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

var (
	// Give me:
	// - All the numbers in a string
	// - That are either integers or decimals of any precision, that are not percentages and are not RR valuess
	numeric = regexp2.MustCompile(`(\b\d+(?:[\.,]\d+)?\b(?!(?:[\.,]\d+)|(?:\s*(?:%|percent|RR|hr|\$|\.|\/))))`, regexp2.None)

	binanceAssetPairs = map[string]bool{}
)

func captureNumbers(content string) ([]float64, error) {
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

		trimmed := strings.ReplaceAll(match.String(), " ", "")
		if trimmed == "" {
			continue
		}

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

func findSide(content string) (tradeengineproto.TRADE_SIDE, bool) {
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
func findTicker(content string) string {
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
			if containsTicker(childContent) != "" {
				return token
			}
		}

		if _, ok := binanceAssetPairs[token]; ok {
			return token
		}
	}

	return ""
}
