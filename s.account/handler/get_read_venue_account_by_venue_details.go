package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// ReadVenueAccountByVenueAccountDetails ...
func (s *AccountService) ReadVenueAccountByVenueAccountDetails(
	ctx context.Context, in *accountproto.ReadVenueAccountByVenueAccountDetailsRequest,
) (*accountproto.ReadVenueAccountByVenueAccountDetailsResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	case in.Venue == tradeengineproto.VENUE_UNREQUIRED:
		return nil, gerrors.BadParam("missing_param.venue", nil)
	case in.RequestContext == "":
		return nil, gerrors.BadParam("missing_param.request_context", nil)
	}

	errParams := map[string]string{
		"user_id":    in.UserId,
		"actor_id":   in.ActorId,
		"venue":      in.Venue.String(),
		"subaccount": in.Subaccount,
	}

	// Validate the request venue credentials.
	if err := validateVenueAccountDetailsRequest(in); err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_venue_account_by_venue_details", errParams)
	}

	// Read.
	venueAccount, err := dao.ReadVenueAccountByVenueAccountDetails(ctx, in.Venue.String(), in.UserId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_venue_account_by_venue_details.dao", errParams)
	}

	// Marshal.
	proto, err := marshaling.VenueAccountDomainToProtoUnmasked(venueAccount)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_venue_account_by_venue_details.marshal_to_proto", errParams)
	}

	return &accountproto.ReadVenueAccountByVenueAccountDetailsResponse{
		VenueAccount: proto,
	}, nil
}

func validateVenueAccountDetailsRequest(in *accountproto.ReadVenueAccountByVenueAccountDetailsRequest) error {
	switch in.RequestContext {
	case accountproto.RequestContextOrderRequest, accountproto.RequestContextUserRequest:
	default:
		return gerrors.FailedPrecondition("invalid_request_context", nil)
	}

	switch in.Venue {
	case tradeengineproto.VENUE_BINANCE, tradeengineproto.VENUE_FTX, tradeengineproto.VENUE_BITFINEX, tradeengineproto.VENUE_DERIBIT:
	default:
		return gerrors.Unimplemented("venue.unimplemented", nil)
	}

	switch in.ActorId {
	case in.UserId:
	case accountproto.ActorSystemTradeEngine, accountproto.ActorSystemPayments:
	default:
		return gerrors.FailedPrecondition("bad_actor", nil)
	}

	return nil
}
