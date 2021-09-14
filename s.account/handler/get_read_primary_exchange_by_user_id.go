package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/domain"
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
	account, err := dao.ReadAccountByUserID(ctx, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "account_not_found"):
		return nil, gerrors.FailedPrecondition("failed_to_read_primary_exchange.account_required", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_read_primary_exchange.dao.read_account", errParams)
	}

	// List exchanges.
	exchanges, err := dao.ListExchangesByUserID(ctx, in.UserId, true)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "exchanges_not_found_for_user_id"):
		return nil, gerrors.Augment(err, "failed_to_read_primary_exchange.no_exchange_found", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_read_primary_exchange.dao.read_primary_exchange", errParams)
	}

	// Parse primary exchange.
	var primaryExchange *domain.Exchange
	for _, exchange := range exchanges {
		if exchange.ExchangeType == account.PrimaryExchange {
			primaryExchange = exchange
		}
	}

	switch {
	case primaryExchange == nil:
		return nil, gerrors.FailedPrecondition("exchange_found_different_to_primary_exchange_on_account", errParams)
	}

	protoExchange, err := marshaling.ExchangeDomainToProto(primaryExchange)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_primary_exchange.marshal_from_domain_to_proto", errParams)
	}

	return &accountproto.ReadPrimaryExchangeByUserIDResponse{
		PrimaryExchange: protoExchange,
	}, nil
}
