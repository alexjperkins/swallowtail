package sync

import (
	"strings"
)

var (
	// Read concurrently; be we don't update, so we avoid data races.
	validAssetPairs = map[string]bool{
		"USDT": true,
		"USD":  true,
		"BTC":  true,
		"ETH":  true,
		"GBP":  true,
	}
)

func isValidAssetPairOrConvert(assetPair string) (string, bool) {
	_, ok := validAssetPairs[assetPair]
	if !ok {
		return "", false
	}
	if strings.ToLower(assetPair) == "usdt" {
		return "usd", true
	}
	return assetPair, true
}
