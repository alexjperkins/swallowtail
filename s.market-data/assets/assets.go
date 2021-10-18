package assets

type AssetVolatiltyRating int

const (
	AssetVolatiltyRatingLow AssetVolatiltyRating = iota
	AssetVolatiltyRatingMedium
	AssetVolatiltyRatingHigh
)

func (a AssetVolatiltyRating) PercentageTriggerValue() float64 {
	switch a {
	case AssetVolatiltyRatingHigh:
		return 0.05
	case AssetVolatiltyRatingMedium:
		return 0.025
	case AssetVolatiltyRatingLow:
		return 0.001
	default:
		return 0.001
	}
}

type AssetPair struct {
	Symbol           string
	AssetPair        string
	VolatilityRating AssetVolatiltyRating
}

var (
	// latestpriceassets are all the assets that we'd like to watch in order to publish
	// information to at a later point in time.
	LatestPriceAssets = []*AssetPair{
		{
			Symbol:           "btc",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingLow,
		},
		{
			Symbol:           "eth",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
		},
		{
			Symbol:           "eth",
			AssetPair:        "btc",
			VolatilityRating: AssetVolatiltyRatingMedium,
		},
		{
			Symbol:           "sol",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
		},
		{
			Symbol:           "avax",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
		},
		{
			Symbol:           "algo",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "cope",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "link",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "srm",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "ray",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "step",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "dot",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
		},
		{
			Symbol:           "aave",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
		},
		{
			Symbol:           "atom",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
		},
		{
			Symbol:           "bnb",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
		},
		{
			Symbol:           "cake",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
		},
		{
			Symbol:           "ftm",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "rune",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "sushi",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "uni",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
		},
		{
			Symbol:           "woo",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "spell",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "dydx",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "luna",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
	}
)
