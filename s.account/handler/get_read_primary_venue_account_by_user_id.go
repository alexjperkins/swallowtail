package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/domain"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
)

// ReadPrimaryVenueAccountByUserID ...
func (s *AccountService) ReadPrimaryVenueAccountByUserID(
	ctx context.Context, in *accountproto.ReadPrimaryVenueAccountByUserIDRequest,
) (*accountproto.ReadPrimaryVenueAccountByUserIDResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	}

	errParams := map[string]string{
		"user_id":  in.UserId,
		"actor_id": in.ActorId,
	}

	// Validate that the user first has an account registered.
	account, err := dao.ReadAccountByUserID(ctx, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "account_not_found"):
		return nil, gerrors.FailedPrecondition("failed_to_read_primary_venue_account.account_required", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_read_primary_venue_account.dao.read_account", errParams)
	}

	// List venue accounts.
	venueAccounts, err := dao.ListVenueAccountsByUserID(ctx, in.UserId, true)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "venue_accounts_not_found_for_user_id"):
		return nil, gerrors.Augment(err, "failed_to_read_primary_venue_account.no_venue_account_found", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_read_primary_venue_account.dao.read_primary_venue_account", errParams)
	}

	// Parse primary venue account.
	var primaryVenueAccount *domain.VenueAccount
	for _, venueAccount := range venueAccounts {
		if venueAccount.VenueID == account.PrimaryVenueAccount {
			primaryVenueAccount = venueAccount
		}
	}

	switch {
	case primaryVenueAccount == nil:
		return nil, gerrors.FailedPrecondition("venue_account_found_different_to_primary_venue_account_on_account", errParams)
	}

	proto, err := marshaling.VenueAccountDomainToProtoUnmasked(primaryVenueAccount)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_primary_venue_account.marshal_from_domain_to_proto_unmasked", errParams)
	}

	return &accountproto.ReadPrimaryVenueAccountByUserIDResponse{
		PrimaryVenueAccount: proto,
	}, nil
}
