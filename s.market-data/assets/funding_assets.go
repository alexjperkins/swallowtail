package assets

import (
	"math"
	accountproto "swallowtail/s.account/proto"
	"sync"
)

// FundingRateAsset ...
type FundingRateAsset struct {
	Symbol   string
	Exchange accountproto.ExchangeType
}

// FundingRateExchangeInfo ...
type FundingRateExchangeInfo struct {
	HigherBound float64
	LowerBound  float64
}

var (
	coeffMu                 sync.RWMutex
	fundingRateExchangeData = map[accountproto.ExchangeType]*FundingRateExchangeInfo{
		accountproto.ExchangeType_BINANCE: {
			HigherBound: 0.4,
			LowerBound:  0.025,
		},
		accountproto.ExchangeType_FTX: {
			HigherBound: 0.01,
			LowerBound:  0.0,
		},
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
func GetFundingRateCoefficientByExchange(exchange accountproto.ExchangeType) *FundingRateExchangeInfo {
	coeffMu.RLock()
	defer coeffMu.RUnlock()

	if v, ok := fundingRateExchangeData[exchange]; ok {
		return v
	}

	return &FundingRateExchangeInfo{
		HigherBound: math.MaxFloat64,
		LowerBound:  -math.MaxFloat64,
	}
}
