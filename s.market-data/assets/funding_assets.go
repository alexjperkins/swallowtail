package assets

import (
	"math"
	"sync"

	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// FundingRateAsset ...
type FundingRateAsset struct {
	Symbol          string
	Venue           tradeengineproto.VENUE
	HumanizedSymbol string
}

// FundingRateVenueInfo ...
type FundingRateVenueInfo struct {
	HigherBound float64
	LowerBound  float64
}

var (
	coeffMu              sync.RWMutex
	fundingRateVenueData = map[tradeengineproto.VENUE]*FundingRateVenueInfo{
		tradeengineproto.VENUE_BINANCE: {
			HigherBound: 0.4,
			LowerBound:  0.025,
		},
		tradeengineproto.VENUE_FTX: {
			HigherBound: 0.01,
			LowerBound:  0.0,
		},
		tradeengineproto.VENUE_BITFINEX: {
			HigherBound: 0.025,
			LowerBound:  -0.005,
		},
	}
)

var (
	// FundingRateAssets ...
	FundingRateAssets = []*FundingRateAsset{
		{
			Symbol: "AVAX-PERP",
			Venue:  tradeengineproto.VENUE_FTX,
		},
		{
			Symbol: "AVAXUSDT",
			Venue:  tradeengineproto.VENUE_BINANCE,
		},
		{
			Symbol: "BTC-PERP",
			Venue:  tradeengineproto.VENUE_FTX,
		},
		{
			Symbol: "BTCUSDT",
			Venue:  tradeengineproto.VENUE_BINANCE,
		},
		{
			Symbol: "ETH-PERP",
			Venue:  tradeengineproto.VENUE_FTX,
		},
		{
			Symbol: "ETHUSDT",
			Venue:  tradeengineproto.VENUE_BINANCE,
		},
		{
			Symbol: "LUNA-PERP",
			Venue:  tradeengineproto.VENUE_FTX,
		},
		{
			Symbol: "LUNAUSDT",
			Venue:  tradeengineproto.VENUE_BINANCE,
		},
		{
			Symbol: "SOL-PERP",
			Venue:  tradeengineproto.VENUE_FTX,
		},
		{
			Symbol: "SOLUSDT",
			Venue:  tradeengineproto.VENUE_BINANCE,
		},
		{
			Symbol: "FTMUSDT",
			Venue:  tradeengineproto.VENUE_BINANCE,
		},
		{
			Symbol: "FTM-PERP",
			Venue:  tradeengineproto.VENUE_FTX,
		},
		{
			Symbol: "ATOMUSDT",
			Venue:  tradeengineproto.VENUE_BINANCE,
		},
		{
			Symbol: "ATOM-PERP",
			Venue:  tradeengineproto.VENUE_FTX,
		},
		{
			Symbol: "ALGOUSDT",
			Venue:  tradeengineproto.VENUE_BINANCE,
		},
		{
			Symbol: "ALGO-PERP",
			Venue:  tradeengineproto.VENUE_FTX,
		},
		{
			Symbol:          "tBTCF0:USTF0",
			Venue:           tradeengineproto.VENUE_BITFINEX,
			HumanizedSymbol: "BTCUSD",
		},
		{
			Symbol:          "tETHF0:USTF0",
			Venue:           tradeengineproto.VENUE_BITFINEX,
			HumanizedSymbol: "ETHUSD",
		},
	}
)

// GetFundingRateCoefficientByVenue ...
func GetFundingRateCoefficientByVenue(venue tradeengineproto.VENUE) *FundingRateVenueInfo {
	coeffMu.RLock()
	defer coeffMu.RUnlock()

	if v, ok := fundingRateVenueData[venue]; ok {
		return v
	}

	return &FundingRateVenueInfo{
		HigherBound: math.MaxFloat64,
		LowerBound:  -math.MaxFloat64,
	}
}
