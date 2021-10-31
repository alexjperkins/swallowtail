package assets

import solananftsproto "swallowtail/s.solana-nfts/proto"

// SolanaNFTInfo ...
type SolanaNFTInfo struct {
	CollectionID          string
	Vendor                solananftsproto.SolanaNFTVendor
	HumanizedCollectionID string
	Emoji                 string
}

var (
	SolanaNFTAssets = []*SolanaNFTInfo{
		{
			CollectionID:          solananftsproto.SolanartCollectionIDGalacticGeckoSpaceGarage,
			Vendor:                solananftsproto.SolanaNFTVendor_SOLANART,
			HumanizedCollectionID: "Galactic Gecko Space Garage",
			Emoji:                 ":lizard:",
		},
	}
)
