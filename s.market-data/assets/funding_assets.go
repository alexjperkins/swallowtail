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
			Symbol:   "SOL-PERP",
			Exchange: accountproto.ExchangeType_FTX,
		},
		{
			Symbol:   "SOLUSDT",
			Exchange: accountproto.ExchangeType_BINANCE,
		},
	}
)
