package marshaling

import (
	"strconv"
	"swallowtail/s.solana-nfts/dto"
)

// SolanartPriceStatisticsDTOToVendorDTO ...
func SolanartPriceStatisticsDTOToVendorDTO(ins *dto.GetSolanartPriceStatisticsByCollectionIDResponse) *dto.GetVendorPriceStatisticsByCollectionIDResponse {
	vendors := make([]*dto.VendorPriceStatistic, 0, len(*ins))
	for _, in := range *ins {
		vendors = append(vendors, &dto.VendorPriceStatistic{
			Name:          in.Name,
			Price:         in.Price,
			ID:            strconv.Itoa(in.ID),
			LastSoldPrice: in.LastSoldPrice,
			IsForSale:     intToBool(in.ForSale),
		})
	}

	return &dto.GetVendorPriceStatisticsByCollectionIDResponse{
		Stats: vendors,
	}
}

func intToBool(v int) bool {
	if v == 1 {
		return true
	}
	return false
}
