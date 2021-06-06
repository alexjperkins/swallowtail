package satoshi

import (
	"math/rand"
)

var (
	insults = []string{
		"Come on, at least give me a ticker such as ETHUSDT.",
		"Mate, that is poggers. Give me a ticker like BTCUSDT",
		"Satoshi didn't build a blockchain for this. Ticker please.",
	}
)

func randomInsultGenerator() string {
	nIndexes := len(insults)
	return insults[rand.Intn(nIndexes)]
}

// calculateRisk returns the number of contracts to buy.
func calculateRisk(entry, stopLoss, accountSize, percentage float64) float64 {
	switch {
	case entry == stopLoss:
		return 0.0
	}
	maxRiskToLose := percentage * accountSize
	lossPerContract := entry - stopLoss
	return maxRiskToLose / lossPerContract
}

func contains(needle string, haystack []string) bool {
	for _, h := range haystack {
		if needle == h {
			return true
		}
	}
	return false
}
