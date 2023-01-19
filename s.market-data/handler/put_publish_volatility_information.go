package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	marketdataproto "swallowtail/s.market-data/proto"
)

// PublishVolatilityInformation ...
func (s *MarketDataService) PublishVolatilityInformation(
	ctx context.Context, in *marketdataproto.PublishVolatilityInformationRequest,
) (*marketdataproto.PublishVolatilityInformationResponse, error) {
	return nil, gerrors.Unimplemented("publish volatility informationk unimplemented", nil)
}
