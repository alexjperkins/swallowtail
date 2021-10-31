package marshaling

import "swallowtail/s.solana-nfts/dto"

// MagicEdenPriceStatisticsDTOToVendorDTO ...
func MagicEdenPriceStatisticsDTOToVendorDTO(in *dto.GetMagicEdenPriceStatisticsByCollectionIDResponse) *dto.GetVendorPriceStatisticsByCollectionIDResponse {
	return &dto.GetVendorPriceStatisticsByCollectionIDResponse{}
}
