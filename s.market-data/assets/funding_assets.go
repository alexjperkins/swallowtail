package assets

import accountproto "swallowtail/s.account/proto"

type FundingRateAsset struct {
	Symbol   string
	Exchange accountproto.ExchangeType
}

var (
	// FundingRateAssets ...
	FundingRateAssets = []*FundingRateAsset{
		{
			Symbol:   "BTCUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "ETHUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "SOLUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "SRMUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "AVAXUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "LUNAUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
		{
			Symbol:   "ATOMUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
	}
)
