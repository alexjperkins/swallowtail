package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
)

// ReadPrimaryExchangeByUserID ...
func (s *AccountService) ReadPrimaryExchangeByUserID(
	ctx context.Context, in *accountproto.ReadPrimaryExchangeByUserIDRequest,
) (*accountproto.ReadPrimaryExchangeByUserIDResponse, error) {
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
	_, err := dao.ReadAccountByUserID(ctx, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "account_not_found"):
		return nil, gerrors.FailedPrecondition("failed_to_read_primary_exchange.account_required", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_read_primary_exchange.dao.read_account", errParams)
	}

	primaryExchange, err := dao.ReadPrimaryExchangeByUserID(ctx, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "not_found.primary_exchange"):
		return nil, gerrors.Augment(err, "failed_to_read_primary_exchange", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_read_primary_exchange.dao.read_primary_exchange", errParams)
	}

	protoExchange, err := marshaling.ExchangeDomainToProto(primaryExchange)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_primary_exchange.marshal_from_domain_to_proto", errParams)
	}

	return &accountproto.ReadPrimaryExchangeByUserIDResponse{
		PrimaryExchange: protoExchange,
	}, nil
}
