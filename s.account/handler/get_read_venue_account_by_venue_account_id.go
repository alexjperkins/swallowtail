package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
)

// ReadVenueAccountByVenueAccountID ...
func (s *AccountService) ReadVenueAccountByVenueAccountID(
	ctx context.Context, in *accountproto.ReadVenueAccountByVenueAccountIDRequest,
) (*accountproto.ReadVenueAccountByVenueAccountIDResponse, error) {
	switch {
	case in.VenueAccountId == "":
		return nil, gerrors.BadParam("missing_param.venue_account_id", nil)
	}

	errParams := map[string]string{
		"venue_account_id": in.VenueAccountId,
	}

	exchange, err := dao.ReadVenueAccountByVenueAccountID(ctx, in.VenueAccountId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_exchange_by_venue_account_id.dao", errParams)
	}

	proto, err := marshaling.VenueAccountDomainToProto(exchange)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_exchange_by_venue_id.marshal_to_proto", errParams)
	}

	return &accountproto.ReadVenueAccountByVenueAccountIDResponse{
		VenueAccount: proto,
	}, nil
}
