package handler

import (
	"context"
	"strconv"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.solana-nfts/client"
	"swallowtail/s.solana-nfts/dto"
	"swallowtail/s.solana-nfts/marshaling"
	solananftsproto "swallowtail/s.solana-nfts/proto"
)

var (
	defaultVendors = []solananftsproto.SolanaNFTVendor{
		solananftsproto.SolanaNFTVendor_MAGIC_EDEN,
		solananftsproto.SolanaNFTVendor_SOLANART,
	}
)

// ReadSolanaPriceStatisticsByCollectionID ...
func (s *SolanaNFTsService) ReadSolanaPriceStatisticsByCollectionID(
	ctx context.Context, in *solananftsproto.ReadSolanaPriceStatisticsByCollectionIDRequest,
) (*solananftsproto.ReadSolanaPriceStatisticsByCollectionIDResponse, error) {
	var vendors = defaultVendors
	switch {
	case in.CollectionId == "":
		return nil, gerrors.BadParam("missing_param.collection_id", nil)
	case in.SearchContext == "":
		return nil, gerrors.BadParam("missing_param.search_context", nil)
	case len(in.FilterByVendors) != 0:
		vendors = in.FilterByVendors
	}

	errParams := map[string]string{
		"collection_id": in.CollectionId,
		"limit":         strconv.Itoa(int(in.Limit)),
		"order":         in.Order.String(),
	}

	// Collect and marshal stats.
	vendorStatsProtos := []*solananftsproto.VendorPriceStatistics{}
	for _, vendor := range vendors {
		rsp, err := client.GetVendorPriceStatisticsByCollectionID(ctx, vendor, &dto.GetVendorPriceStatisticsByCollectionIDRequest{
			CollectionID: in.CollectionId,
		}, in.Order, int(in.Limit))
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_get_vendor_statistics", errParams)
		}

		vendorStatsProto := marshaling.VendorPriceStatisticsDTOToProto(rsp)

		vendorStatsProtos = append(vendorStatsProtos, vendorStatsProto)
	}

	return &solananftsproto.ReadSolanaPriceStatisticsByCollectionIDResponse{
		VendorStats: vendorStatsProtos,
	}, nil
}
