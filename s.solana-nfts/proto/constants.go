package solananftsproto

const (
	SearchContextMarketData = "search-context-market-data"
	SearchContextAdHoc      = "search-context-ad-hoc"
)

const (
	SolanartCollectionIDGalacticGeckoSpaceGarage = "galacticgeckospacegarage"
)

const (
	MagicEndCollectionIDGloomPunks = "gloom_punk_club"
)

// IsValidCollectionIDByVendor defines if the collection ID & the vendor are a valid pair.
func IsValidCollectionIDByVendor(vendor SolanaNFTVendor, collectionID string) bool {
	switch vendor {
	case SolanaNFTVendor_MAGIC_EDEN:
		switch collectionID {
		case MagicEndCollectionIDGloomPunks:
			return true
		default:
			return false
		}
	case SolanaNFTVendor_SOLANART:
		switch collectionID {
		case SolanartCollectionIDGalacticGeckoSpaceGarage:
			return true
		default:
			return false
		}
	default:
		return false
	}
}
