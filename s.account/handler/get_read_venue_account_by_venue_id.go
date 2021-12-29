package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// ReadVenueAccountByVenueID ...
func (s *AccountService) ReadVenueAccountByVenueID(
	ctx context.Context, in *accountproto.ReadVenueAccountByVenueIDRequest,
) (*accountproto.ReadVenueAccountByVenueIDResponse, error) {
	if in.Venue == tradeengineproto.VENUE_UNREQUIRED {
		return nil, gerrors.BadParam("missing_param.venue", nil)
	}

	errParams := map[string]string{
		"venue": in.Venue.String(),
	}

	exchange, err := dao.ReadVenueAccountByVenueID(ctx, in.Venue.String())
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_exchange_by_venue_id.dao", errParams)
	}

	proto, err := marshaling.VenueAccountDomainToProto(exchange)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_exchange_by_venue_id.marshal_to_proto", errParams)
	}

	return &accountproto.ReadVenueAccountByVenueIDResponse{
		VenueAccount: proto,
	}, nil
}
