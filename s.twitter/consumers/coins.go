package consumers

import "sync"

var (
	coinSymbols    = map[string]bool{}
	coinSymbolsMtx sync.RWMutex
)

func getCoinSymbols() []string {
	coinSymbolsMtx.RLock()
	defer coinSymbolsMtx.RUnlock()
	coins := []string{}
	for c := range coinSymbols {
		coins = append(coins, c)
	}
	return coins
}

func addNewCoin(symbol string) bool {
	coinSymbolsMtx.RLock()
	defer coinSymbolsMtx.RUnlock()
	if _, ok := coinSymbols[symbol]; ok {
		return false
	}
	coinSymbols[symbol] = true
	return true
}
