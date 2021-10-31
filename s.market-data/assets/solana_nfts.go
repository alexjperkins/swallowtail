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
		{
			CollectionID:          solananftsproto.SolanartCollectionIDDegenerateApeAcademy,
			Vendor:                solananftsproto.SolanaNFTVendor_SOLANART,
			HumanizedCollectionID: "Degenerate Ape Academy",
			Emoji:                 ":monkey:",
		},
		{
			CollectionID:          solananftsproto.SolanartCollectionIDGloomPunk,
			Vendor:                solananftsproto.SolanaNFTVendor_SOLANART,
			HumanizedCollectionID: "Gloom Punk",
			Emoji:                 ":woman_artist:",
		},
		{
			CollectionID:          solananftsproto.SolanartCollectionIDBabyApes,
			Vendor:                solananftsproto.SolanaNFTVendor_SOLANART,
			HumanizedCollectionID: "Baby Apes",
			Emoji:                 ":monkey_face:",
		},
		{
			CollectionID:          solananftsproto.SolanartCollectionIDSolarmy2D,
			Vendor:                solananftsproto.SolanaNFTVendor_SOLANART,
			HumanizedCollectionID: "Solarmy 2D",
			Emoji:                 ":poop:",
		},
		{
			CollectionID:          solananftsproto.SolanartCollectionIDSolarmy3D,
			Vendor:                solananftsproto.SolanaNFTVendor_SOLANART,
			HumanizedCollectionID: "Solarmy 3D",
			Emoji:                 ":mechanical_arm:",
		},
		{
			CollectionID:          solananftsproto.SolanartCollectionIDThugBirdz,
			Vendor:                solananftsproto.SolanaNFTVendor_SOLANART,
			HumanizedCollectionID: "ThugBirdz",
			Emoji:                 ":bird:",
		},
		{
			CollectionID:          solananftsproto.SolanartCollectionIDFrakt,
			Vendor:                solananftsproto.SolanaNFTVendor_SOLANART,
			HumanizedCollectionID: "Frakt",
			Emoji:                 ":art:",
		},
		{
			CollectionID:          solananftsproto.SolanartCollectionIDTurtles,
			Vendor:                solananftsproto.SolanaNFTVendor_SOLANART,
			HumanizedCollectionID: "Turtles",
			Emoji:                 ":turtle:",
		},
		{
			CollectionID:          solananftsproto.SolanartCollectionIDTheTower,
			Vendor:                solananftsproto.SolanaNFTVendor_SOLANART,
			HumanizedCollectionID: "The Tower DAO",
			Emoji:                 ":tokyo_tower:",
		},
	}
)
