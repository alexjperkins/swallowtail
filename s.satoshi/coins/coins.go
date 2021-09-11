package coins

import "sync"

var (
	coinSymbols = map[string]bool{
		"AVAX":  true,
		"FTM":   true,
		"BTC":   true,
		"ETH":   true,
		"ROPE":  true,
		"LTC":   true,
		"OCEAN": true,
		"RSR":   true,
		"NOIA":  true,
		"HTR":   true,
		"SOL":   true,
		"AKT":   true,
		"BNB":   true,
		"ALPHA": true,
		"WOO":   true,
		"ALGO":  true,
		"AAVE":  true,
		"RUNE":  true,
		"SAND":  true,
		"FET":   true,
		"FTT":   true,
		"RAY":   true,
		"API3":  true,
		"UNI":   true,
		"1INCH": true,
		"BAND":  true,
		"BAL":   true,
		"CAKE":  true,
		"SRM":   true,
		"ORK":   true,
		"AKRO":  true,
		"TVK":   true,
		"BOSON": true,
		"FIDA":  true,
		"OXY":   true,
		"YFI":   true,
		"MIR":   true,
		"CRWNY": true,
		"STEP":  true,
		"LINK":  true,
		"MEDIA": true,
		"LIQ":   true,
		"FAB":   true,
		"COPE":  true,
	}
	coinSymbolsMtx sync.RWMutex
)

// List gets all coin symbols for satoshi.
func List() []string {
	coinSymbolsMtx.RLock()
	defer coinSymbolsMtx.RUnlock()
	coins := []string{}
	for c := range coinSymbols {
		coins = append(coins, c)
	}
	return coins
}

// AddNewCoin adds a new coin.
func AddNewCoin(symbol string) bool {
	coinSymbolsMtx.RLock()
	defer coinSymbolsMtx.RUnlock()
	if _, ok := coinSymbols[symbol]; ok {
		return false
	}
	coinSymbols[symbol] = true
	return true
}
