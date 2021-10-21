package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	ftxproto "swallowtail/s.ftx/proto"
)

// GetFundingRate ...
func (s *FTXService) GetFundingRate(
	ctx context.Context,
	in *ftxproto.GetFTXFundingRatesRequest,
) (*ftxproto.GetFTXFundingRatesResponse, error) {
	return nil, gerrors.Unimplemented("failed_to_get_funding_rate.unimplemented", nil)
}
