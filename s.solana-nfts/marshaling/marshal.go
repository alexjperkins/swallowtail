package marshaling

import (
	"swallowtail/s.solana-nfts/dto"
	solananftsproto "swallowtail/s.solana-nfts/proto"
)

// PriceStatisticDTOToProtos ...
func PriceStatisticDTOToProtos(in *dto.GetVendorPriceStatisticsByCollectionIDResponse, vendor solananftsproto.SolanaNFTVendor) []*solananftsproto.PriceStatistic {
	var protos = make([]*solananftsproto.PriceStatistic, 0, len(in.Stats))
	for _, stat := range in.Stats {
		protos = append(protos, &solananftsproto.PriceStatistic{
			Price:         float32(stat.Price),
			Id:            stat.ID,
			LastSoldPrice: float32(stat.LastSoldPrice),
			Name:          stat.Name,
			ForSale:       stat.IsForSale,
			Vendor:        vendor,
		})
	}

	return protos
}
