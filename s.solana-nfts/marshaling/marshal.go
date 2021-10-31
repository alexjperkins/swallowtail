package marshaling

import (
	"swallowtail/s.solana-nfts/dto"
	solananftsproto "swallowtail/s.solana-nfts/proto"
)

// VendorPriceStatisticsDTOToProto ...
func VendorPriceStatisticsDTOToProto(in *dto.GetVendorPriceStatisticsByCollectionIDResponse) *solananftsproto.VendorPriceStatistics {
	return &solananftsproto.VendorPriceStatistics{}
}
