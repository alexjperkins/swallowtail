package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// ReadInternalVenueAccount ...
func (s *AccountService) ReadInternalVenueAccount(
	ctx context.Context, in *accountproto.ReadInternalVenueAccountRequest,
) (*accountproto.ReadInternalVenueAccountResponse, error) {
	switch {
	case in.Venue == tradeengineproto.VENUE_UNREQUIRED:
		return nil, gerrors.BadParam("missing_param.venue", nil)
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	}

	if ok := isValidActorID(in.ActorId); !ok {
		return nil, gerrors.Unauthenticated("failed_to_read_internal_venue_account.unauthenticated", nil)
	}

	errParams := map[string]string{
		"actor_id":           in.ActorId,
		"venue":              in.Venue.String(),
		"subaccount":         in.Subaccount,
		"venue_account_type": in.VenueAccountType.String(),
	}

	internalVenueAccount, err := dao.ReadInternalVenueAccount(ctx, in.Venue.String(), in.Subaccount, in.VenueAccountType.String())
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_internal_venue_account.dao", errParams)
	}

	proto, err := marshaling.InternalVenueAccountDomainToProto(internalVenueAccount)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_internal_venue_account.marshal", errParams)
	}

	return &accountproto.ReadInternalVenueAccountResponse{
		InternalVenueAccount: proto,
	}, nil
}
