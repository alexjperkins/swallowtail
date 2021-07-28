package handler

import (
	"context"

	"github.com/monzo/terrors"

	"swallowtail/s.market-data/clients/bybt"
	"swallowtail/s.market-data/marshaling"
	marketdataproto "swallowtail/s.market-data/proto"
)

// GetExchangeFundingRates ...
func (s *MarketDataService) GetExchangeFundingRates(
	ctx context.Context, in *marketdataproto.GetExchangeFundingRatesRequest,
) (*marketdataproto.GetExchangeFundingRatesResponse, error) {
	switch {
	case in.Asset == "":
		return nil, terrors.PreconditionFailed("empty-asset", "Cannot get funding rates for null asset", nil)
	}

	errParams := map[string]string{
		"asset": in.GetAsset(),
	}

	rsp, err := bybt.GetExchangeFundingRatesByAsset(ctx, &bybt.GetExchangeFundingRatesByAssetRequest{
		Asset: in.Asset,
	})
	if err != nil {
		return nil, terrors.Augment(err, "Failed to get exchanges funding rates", errParams)
	}

	return marshaling.GetExchangeFundingRatesProtoToDTO(rsp), nil
}
