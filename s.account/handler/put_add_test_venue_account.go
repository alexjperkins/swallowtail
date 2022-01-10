package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
)

// AddTestVenueAccount ...
func (s *AccountService) AddTestVenueAccount(
	ctx context.Context, in *accountproto.AddTestVenueAccountRequest,
) (*accountproto.AddTestVenueAccountResponse, error) {
	// Validation.
	switch {
	case in.VenueAccount == nil:
		return nil, gerrors.BadParam("missing_param.venue_account", nil)
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	}

	errParams := map[string]string{
		"user_id":  in.UserId,
		"actor_id": in.ActorId,
	}

	// Validate actor.
	if err := validateAccountCreationActor(in.ActorId); err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_test_venue_account.unauthorized", errParams)
	}

	// Marshal to domain.
	domainVenueAccount, err := marshaling.VenueAccountProtoToDomain(in.UserId, in.VenueAccount)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_test_venue_account.marshaling", errParams)
	}

	// TODO: validate account.

	// Create or update test account.
	venueAccount, err := dao.CreateOrUpdateTestVenueAccount(ctx, domainVenueAccount, in.PreventOverwrite)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_test_venue_account.dao", errParams)
	}

	proto, err := marshaling.VenueAccountDomainToProto(venueAccount)

	return &accountproto.AddTestVenueAccountResponse{
		VenueAccount: proto,
	}, nil

}

// TODO: move to `helpers.go`
func validateAccountCreationActor(actorID string) error {
	return nil
}
