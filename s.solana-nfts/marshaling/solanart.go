package marshaling

import "swallowtail/s.solana-nfts/dto"

// SolanartPriceStatisticsDTOToVendorDTO ...
func SolanartPriceStatisticsDTOToVendorDTO(in *dto.GetSolanartPriceStatisticsByCollectionIDResponse) *dto.GetVendorPriceStatisticsByCollectionIDResponse {
	return &dto.GetVendorPriceStatisticsByCollectionIDResponse{}
}
