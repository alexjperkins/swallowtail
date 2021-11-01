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

// AssetPair defines an asset.
type AssetPair struct {
	Symbol           string
	AssetPair        string
	VolatilityRating AssetVolatiltyRating
	Grouping         AssetGroup
}

var (
	// latestpriceassets are all the assets that we'd like to watch in order to publish
	// information to at a later point in time.
	LatestPriceAssets = []*AssetPair{
		{
			Symbol:           "btc",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingLow,
			Grouping:         AssetGroupBitcoin,
		},
		{
			Symbol:           "eth",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupEthereum,
		},
		{
			Symbol:           "eth",
			AssetPair:        "btc",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupEthereum,
		},
		{
			Symbol:           "sol",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupSolanaEcosystem,
		},
		{
			Symbol:           "avax",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupLayerOne,
		},
		{
			Symbol:           "algo",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupLayerOne,
		},
		{
			Symbol:           "cope",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupSolanaEcosystem,
		},
		{
			Symbol:           "link",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupOracles,
		},
		{
			Symbol:           "srm",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupSolanaEcosystem,
		},
		{
			Symbol:           "ray",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupSolanaEcosystem,
		},
		{
			Symbol:           "step",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupSolanaEcosystem,
		},
		{
			Symbol:           "dot",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupLayerOne,
		},
		{
			Symbol:           "aave",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupDeFi,
		},
		{
			Symbol:           "atom",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupLayerOne,
		},
		{
			Symbol:           "bnb",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupBSC,
		},
		{
			Symbol:           "cake",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupBSC,
		},
		{
			Symbol:           "ftm",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupLayerOne,
		},
		{
			Symbol:           "rune",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupDeFi,
		},
		{
			Symbol:           "sushi",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupDeFi,
		},
		{
			Symbol:           "uni",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupDeFi,
		},
		{
			Symbol:           "woo",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupDeFi,
		},
		{
			Symbol:           "spell",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupDeFi2,
		},
		{
			Symbol:           "dydx",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupDeFi2,
		},
		{
			Symbol:           "luna",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupLuna,
		},
		{
			Symbol:           "liq",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupSolanaEcosystem,
		},
		{
			Symbol:           "fab",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupSolanaEcosystem,
		},
		{
			Symbol:           "bop",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupSolanaEcosystem,
		},
		{
			Symbol:           "sol",
			AssetPair:        "btc",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupSolanaEcosystem,
		},
		{
			Symbol:           "sol",
			AssetPair:        "eth",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupSolanaEcosystem,
		},
		{
			Symbol:           "axs",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupMetaverse,
		},
		{
			Symbol:           "doge",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupDogcoin,
		},
		{
			Symbol:           "samo",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupDogcoin,
		},
		{
			Symbol:           "floki",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupDogcoin,
		},
		{
			Symbol:           "scrt",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupPrivacy,
		},
		{
			Symbol:           "one",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupLayerOne,
		},
		{
			Symbol:           "frkt",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
		},
		{
			Symbol:           "htr",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupLayerOne,
		},
		{
			Symbol:           "crv",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupDeFi,
		},
		{
			Symbol:           "comp",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingMedium,
			Grouping:         AssetGroupDeFi,
		},
		{
			Symbol:           "cheems",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupDogcoin,
		},
		{
			Symbol:           "mana",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupMetaverse,
		},
		{
			Symbol:           "sand",
			AssetPair:        "usd",
			VolatilityRating: AssetVolatiltyRatingHigh,
			Grouping:         AssetGroupMetaverse,
		},
	}
)
