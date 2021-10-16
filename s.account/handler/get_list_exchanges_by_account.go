package handler

import (
	"context"
	"strconv"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/domain"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
)

// ListExchanges ...
func (s *AccountService) ListExchanges(
	ctx context.Context, in *accountproto.ListExchangesRequest,
) (*accountproto.ListExchangesResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.FailedPrecondition("missing-param.user_id", nil)
	case !isValidActorUnmaskedRequest(in.ActorId, in.WithUnmaaskedCredentials):
		return nil, gerrors.Unauthenticated("failed_to_list_exchanges_by_user_id.unauthenticated", map[string]string{
			"user_id":  "user_id",
			"actor_id": "actor_id",
		})
	}

	errParams := map[string]string{
		"user_id":  in.UserId,
		"actor_id": in.ActorId,
		"unmasked": strconv.FormatBool(in.WithUnmaaskedCredentials),
	}

	// Determine the correct marshaller to user depending on whether the requester is allowed
	// unmasked credentials.
	var marshaller func([]*domain.Exchange) ([]*accountproto.Exchange, error)
	switch {
	case in.WithUnmaaskedCredentials:
		marshaller = marshaling.ExchangeDomainsToProtosUnmasked
	default:
		marshaller = marshaling.ExchangeDomainsToProtos
	}

	exchanges, err := dao.ListExchangesByUserID(ctx, in.UserId, in.GetActiveOnly())
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "exchanges_not_found_for_user_id"):
		return nil, gerrors.Augment(err, "failed_to_list_exchanges_by_user_id", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_read_exchanges_by_user_id", errParams)
	}

	protos, err := marshaller(exchanges)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_exchanges.at_least_one_exchange_has_an_unsupported_exchange_type", errParams)
	}

	return &accountproto.ListExchangesResponse{
		Exchanges: protos,
	}, nil
}
