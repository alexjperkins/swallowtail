package handler

import (
	"context"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
)

// AddVenueAccount ...
func (s *AccountService) AddVenueAccount(
	ctx context.Context, in *accountproto.AddVenueAccountRequest,
) (*accountproto.AddVenueAccountResponse, error) {
	// Validation.
	switch {
	case in.GetVenueAccount() == nil:
		return nil, gerrors.BadParam("missing_param.venue_account", nil)
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	}

	// Validate venue account.
	if err := validateVenueAccount(in.VenueAccount); err != nil {
		return nil, gerrors.Augment(err, "failed to validate venue account", nil)
	}

	errParams := map[string]string{
		"user_id":    in.UserId,
		"venue":      in.VenueAccount.Venue.String(),
		"subaccount": in.VenueAccount.SubAccount,
	}

	// Confirm the requester first has an account with us.
	_, err := dao.ReadAccountByUserID(ctx, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "account_not_found"):
		return nil, gerrors.FailedPrecondition("cannot_add_venue_account_before_account_created", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "add_venue_account_request_failed.failed_to_read_account_by_user_id", errParams)
	}

	// Check the user hasn't already reached the maximum number of registered venue accounts.
	exs, err := dao.ListVenueAccountsByUserID(ctx, in.UserId, true)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "venue_accounts_not_found_for_user_id"):
	case err != nil:
		return nil, gerrors.Augment(err, "add_venue_account_request_failed.failed_read_existing_registered_venue_account_by_user_id", errParams)
	case len(exs) >= 10:
		return nil, gerrors.FailedPrecondition("add_venue_account_request_failed.maximum_regsitered_active_venue_accounts_reached", errParams)
	}

	// Verify the credentials actually work before storing them in persistent storage.
	verified, reason, err := validateVenueCredentials(ctx, in.UserId, in.VenueAccount)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_request.venue_account_credentials_validation", errParams)
	}

	if !verified {
		errParams["reason"] = reason
		slog.Warn(ctx, "Failed to verify users venue account credentials for %s: %s: invalid", in.VenueAccount.Venue, in.UserId)
		return &accountproto.AddVenueAccountResponse{
			Verified: false,
			Reason:   reason,
		}, nil
	}

	// Marshal to domain.
	venueAccount, err := marshaling.VenueAccountProtoToDomain(in.UserId, in.VenueAccount)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_marshal_request", errParams)
	}

	// Persist venue account.
	if err := dao.AddVenueAccount(ctx, venueAccount); err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_venue_account_for_user", errParams)
	}

	slog.Info(ctx, "Added new venue account for user, with verified credentials", errParams)

	// Mask keys before returning.
	in.VenueAccount.ApiKey = util.MaskKey(in.VenueAccount.ApiKey, 4)
	in.VenueAccount.SecretKey = util.MaskKey(in.VenueAccount.SecretKey, 4)

	return &accountproto.AddVenueAccountResponse{
		VenueAccount: in.VenueAccount,
		Verified:     true,
		// Passing the reason even if verified; since there are some cases where we want to validate the credentials, but also pass a warning message.
		Reason: reason,
	}, nil
}
