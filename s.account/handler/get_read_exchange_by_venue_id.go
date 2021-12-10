package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// ReadExchangeByVenueID ...
func (s *AccountService) ReadExchangeByVenueID(
	ctx context.Context, in *accountproto.ReadExchangeByVenueIDRequest,
) (*accountproto.ReadExchangeByVenueIDResponse, error) {
	if in.Venue == tradeengineproto.VENUE_UNREQUIRED {
		return nil, gerrors.BadParam("missing_param.venue_id", nil)
	}

	errParams := map[string]string{
		"venue": in.Venue.String(),
	}

	exchange, err := dao.ReadExchangeByVenueID(ctx, in.Venue.String())
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_exchange_by_venue_id.dao", errParams)
	}

	proto, err := marshaling.ExchangeDomainToProto(exchange)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_exchange_by_venue_id.marshal_to_proto", errParams)
	}

	return &accountproto.ReadExchangeByVenueIDResponse{
		Exchange: proto,
	}, nil
}
