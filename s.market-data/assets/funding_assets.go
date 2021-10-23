package assets

import accountproto "swallowtail/s.account/proto"

// FundingRateAsset ...
type FundingRateAsset struct {
	Symbol   string
	Exchange accountproto.ExchangeType
}

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
