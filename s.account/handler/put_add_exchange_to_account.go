package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

// AddExchange ...
func (s *AccountService) AddExchange(
	ctx context.Context, in *accountproto.AddExchangeRequest,
) (*accountproto.AddExchangeResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.GetExchange() == nil:
		return nil, gerrors.BadParam("missing_param.exchange", nil)
	case in.Exchange.ApiKey == "":
		return nil, gerrors.BadParam("missing_param.api_key", nil)
	case in.Exchange.SecretKey == "":
		return nil, gerrors.BadParam("missing_param.secret_key", nil)
	case in.Exchange.ExchangeType == accountproto.ExchangeType_FTX:
		return nil, gerrors.Unimplemented("ftx_exchange_unimplemented.coming_shortly", nil)
	}

	errParams := map[string]string{
		"user_id":       in.UserId,
		"exchange_type": in.Exchange.ExchangeType.String(),
	}

	// Confirm the requester first has an account with us.
	_, err := dao.ReadAccountByUserID(ctx, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "account_not_found"):
		return nil, gerrors.FailedPrecondition("cannot_add_exchange_information_before_account_created", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "add_exchange_request_failed.failed_to_read_account_by_user_id", errParams)
	}

	// Check the user hasn't already reached the maximum number of registered exchanges.
	exs, err := dao.ListExchangesByUserID(ctx, in.UserId, true)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "exchanges_not_found_for_user_id"):
	case err != nil:
		return nil, gerrors.Augment(err, "add_exchange_request_failed.failed_read_existing_registered_exchanges_by_user_id", errParams)
	case len(exs) >= 5:
		return nil, gerrors.FailedPrecondition("add_exchange_request_failed.maximum_regsitered_active_exchanges_reached", errParams)
	}

	exchange, err := marshaling.ExchangeProtoToDomain(in.Exchange)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to marshal request", errParams)
	}

	// Verify the credentials actually work before storing them in persistent storage.
	verified, reason, err := validateExchangeCredentials(ctx, exchange)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_request.exchange_validation", errParams)
	}

	if !verified {
		return &accountproto.AddExchangeResponse{
			Verified: false,
			Reason:   reason,
		}, nil
	}

	if err := dao.AddExchange(ctx, exchange); err != nil {
		return nil, terrors.Augment(err, "Failed to add exchange to account.", errParams)
	}

	slog.Info(ctx, "Added new exchange to account", errParams)

	return &accountproto.AddExchangeResponse{
		Exchange: in.Exchange,
		Verified: true,
		// Passing the reason even if verified; since there are some cases where we want to validate the credentials, but also pass a warning message.
		Reason: reason,
	}, nil
}
