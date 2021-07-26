package handler

import (
	"context"

	"github.com/monzo/terrors"

	"swallowtail/s.bybt/client"
	"swallowtail/s.bybt/marshaling"
	bybtproto "swallowtail/s.bybt/proto"
)

// GetExchangeFundingRates ...
func (s *ByBtService) GetExchangeFundingRates(
	ctx context.Context, in *bybtproto.GetExchangeFundingRatesRequest,
) (*bybtproto.GetExchangeFundingRatesResponse, error) {
	switch {
	case in.Asset == "":
		return nil, terrors.PreconditionFailed("empty-asset", "Cannot get funding rates for null asset", nil)
	}

	errParams := map[string]string{
		"asset": in.GetAsset(),
	}

	rsp, err := client.GetExchangeFundingRatesByAsset(ctx, &client.GetExchangeFundingRatesByAssetRequest{
		Asset: in.Asset,
	})
	if err != nil {
		return nil, terrors.Augment(err, "Failed to get exchanges funding rates", errParams)
	}

	return marshaling.GetExchangeFundingRatesProtoToDTO(rsp), nil
}
