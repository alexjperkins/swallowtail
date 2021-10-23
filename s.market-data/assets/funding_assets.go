package assets

import (
	accountproto "swallowtail/s.account/proto"
	"sync"
)

// FundingRateAsset ...
type FundingRateAsset struct {
	Symbol   string
	Exchange accountproto.ExchangeType
}

var (
	coeffMu                         sync.RWMutex
	fundingRateExchangeCoefficients = map[accountproto.ExchangeType]float64{
		accountproto.ExchangeType_BINANCE: 0.1,
		accountproto.ExchangeType_FTX:     1.0,
	}
)

var (
	// FundingRateAssets ...
	FundingRateAssets = []*FundingRateAsset{
		{
			Symbol:   "AVAX-PERP",
			Exchange: accountproto.ExchangeType_FTX,
		},
		{
			Symbol:   "AVAXUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "BTC-PERP",
			Exchange: accountproto.ExchangeType_FTX,
		},
		{
			Symbol:   "BTCUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "ETH-PERP",
			Exchange: accountproto.ExchangeType_FTX,
		},
		{
			Symbol:   "ETHUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "LUNA-PERP",
			Exchange: accountproto.ExchangeType_FTX,
		},
		{
			Symbol:   "LUNAUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "SOL-PERP",
			Exchange: accountproto.ExchangeType_FTX,
		},
		{
			Symbol:   "SOLUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
	}
)

// GetFundingRateCoefficientByExchange ...
func GetFundingRateCoefficientByExchange(exchange accountproto.ExchangeType) float64 {
	coeffMu.RLock()
	defer coeffMu.RUnlock()

	if v, ok := fundingRateExchangeCoefficients[exchange]; ok {
		return v
	}

	return 1.0
}
