package assets

import (
	"math"
	accountproto "swallowtail/s.account/proto"
	"sync"
)

// FundingRateAsset ...
type FundingRateAsset struct {
	Symbol          string
	Exchange        accountproto.ExchangeType
	HumanizedSymbol string
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
		accountproto.ExchangeType_BITFINEX: {
			HigherBound: 0.01,
			LowerBound:  -0.01,
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
		{
			Symbol:   "FTMUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "FTM-PERP",
			Exchange: accountproto.ExchangeType_FTX,
		},
		{
			Symbol:   "ATOMUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "ATOM-PERP",
			Exchange: accountproto.ExchangeType_FTX,
		},
		{
			Symbol:   "ALGOUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "ALGO-PERP",
			Exchange: accountproto.ExchangeType_FTX,
		},
		{
			Symbol:          "tBTCF0:USTF0",
			Exchange:        accountproto.ExchangeType_BITFINEX,
			HumanizedSymbol: "BTCUSD",
		},
		{
			Symbol:          "tETHF0:USTF0",
			Exchange:        accountproto.ExchangeType_BITFINEX,
			HumanizedSymbol: "ETHUSD",
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
