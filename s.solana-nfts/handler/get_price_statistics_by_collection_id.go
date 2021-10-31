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

// ReadSolanaPriceStatisticsByCollectionID ...
func (s *SolanaNFTsService) ReadSolanaPriceStatisticsByCollectionID(
	ctx context.Context, in *solananftsproto.ReadSolanaPriceStatisticsByCollectionIDRequest,
) (*solananftsproto.ReadSolanaPriceStatisticsByCollectionIDResponse, error) {
	// Validation.
	switch {
	case in.CollectionId == "":
		return nil, gerrors.BadParam("missing_param.collection_id", nil)
	case in.SearchContext == "":
		return nil, gerrors.BadParam("missing_param.search_context", nil)
	case in.Vendor == solananftsproto.SolanaNFTVendor_UNKNOWN:
		return nil, gerrors.BadParam("missing_param.vendor", nil)
	case !solananftsproto.IsValidCollectionIDByVendor(in.Vendor, in.CollectionId):
		return nil, gerrors.BadParam("bad_param.collection_id.not_valid_for_vendor", map[string]string{
			"collection_id": in.CollectionId,
			"vendor":        in.Vendor.String(),
		})
	}

	errParams := map[string]string{
		"collection_id": in.CollectionId,
		"limit":         strconv.Itoa(int(in.Limit)),
		"order":         in.Order.String(),
		"vendor":        in.Vendor.String(),
	}

	// Collect stats.
	rsp, err := client.GetVendorPriceStatisticsByCollectionID(ctx, in.Vendor, &dto.GetVendorPriceStatisticsByCollectionIDRequest{
		CollectionID: in.CollectionId,
	}, in.Order, int(in.Limit))
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_vendor_statistics", errParams)
	}

	// Marshal.
	protos := marshaling.PriceStatisticDTOToProtos(rsp)
	return &solananftsproto.ReadSolanaPriceStatisticsByCollectionIDResponse{
		VendorStats: protos,
	}, nil
}
