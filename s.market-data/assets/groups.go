package assets

type AssetGroup int

const (
	AssetGroupOther AssetGroup = iota
	AssetGroupMetaverse
	AssetGroupDeFi
	AssetGroupDeFi2
	AssetGroupSolanaEcosystem
	AssetGroupLayerOne
	AssetGroupBitcoin
	AssetGroupEthereum
	AssetGroupLuna
	AssetGroupBSC
	AssetGroupEthereumLayerTwo
	AssetGroupOracles
	AssetGroupDogcoin
	AssetGroupPrivacy
)

func (a AssetGroup) String() string {
	switch a {
	case AssetGroupMetaverse:
		return "metaverse"
	case AssetGroupDeFi:
		return "defi"
	case AssetGroupDeFi2:
		return "defi2.0"
	case AssetGroupSolanaEcosystem:
		return "solana"
	case AssetGroupLayerOne:
		return "L1"
	case AssetGroupBitcoin:
		return "bitcoin"
	case AssetGroupEthereum:
		return "ethereum"
	case AssetGroupBSC:
		return "bsc"
	case AssetGroupEthereumLayerTwo:
		return "ethereum-layer2"
	case AssetGroupOracles:
		return "oracles"
	case AssetGroupDogcoin:
		return "dogcoin"
	case AssetGroupPrivacy:
		return "privacy"
	case AssetGroupLuna:
		return "luna"
	default:
		return "other"
	}
}
