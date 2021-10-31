package marshaling

import (
	"swallowtail/s.solana-nfts/dto"
	solananftsproto "swallowtail/s.solana-nfts/proto"
)

// PriceStatisticDTOToProtos ...
func PriceStatisticDTOToProtos(in *dto.GetVendorPriceStatisticsByCollectionIDResponse) []*solananftsproto.PriceStatistic {
	var protos []*solananftsproto.PriceStatistic
	return protos
}
