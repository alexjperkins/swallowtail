package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"

	"github.com/monzo/slog"
)

// CreateOrUpdateInternalVenueAccountRequest ...
func (s *AccountService) CreateOrUpdateInternalVenueAccountRequest(
	ctx context.Context, in *accountproto.CreateOrUpdateInternalVenueAccountRequest,
) (*accountproto.CreateOrUpdateInternalVenueAccountResponse, error) {
	// Validation.
	switch {
	case in.InternalVenueAccount == nil:
		return nil, gerrors.BadParam("missing_param.internal_venue_account", nil)
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	}

	// Validate internal venue account.
	if err := validateInternalVenueAccount(in.GetInternalVenueAccount()); err != nil {
		return nil, gerrors.Augment(err, "invalid_internal_venue_account", nil)
	}

	errParams := map[string]string{
		"user_id":  in.UserId,
		"actor_id": in.ActorId,
	}

	// Validate actor.
	if ok := isValidActorID(in.ActorId); !ok {
		return nil, gerrors.FailedPrecondition("failed_to_add_internal_venue_account.unauthorized", errParams)
	}

	// Verify the credentials actually work before storing them in persistent storage.
	verified, reason, err := validateVenueCredentials(ctx, in.UserId, in.InternalVenueAccount)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_request.venue_account_credentials_validation", errParams)
	}

	if !verified {
		errParams["reason"] = reason
		slog.Warn(ctx, "Failed to verify users venue account credentials for %s: %s: invalid", in.InternalVenueAccount.Venue, in.ActorId)
		return &accountproto.CreateOrUpdateInternalVenueAccountResponse{
			Verified: false,
			Reason:   reason,
		}, nil
	}

	// Marshal to domain.
	domainInternalVenueAccount, err := marshaling.InternalVenueAccountProtoToDomain(in.InternalVenueAccount)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_internal_venue_account.marshaling", errParams)
	}

	// Create or update internal account.
	if err := dao.CreateOrUpdateInternalVenueAccount(ctx, domainInternalVenueAccount, in.AllowUpdates); err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_internal_venue_account.dao", errParams)
	}

	slog.Info(ctx, "Created or updated new venue account, with verified credentials")

	// Mask keys before returning.
	in.InternalVenueAccount.ApiKey = util.MaskKey(in.InternalVenueAccount.ApiKey, 4)
	in.InternalVenueAccount.SecretKey = util.MaskKey(in.InternalVenueAccount.SecretKey, 4)

	return &accountproto.CreateOrUpdateInternalVenueAccountResponse{
		InternalVenueAccount: in.InternalVenueAccount,
		Verified:             true,
	}, nil

}
