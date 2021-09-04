package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
)

// ListExchanges ...
func (s *AccountService) ListExchanges(
	ctx context.Context, in *accountproto.ListExchangesRequest,
) (*accountproto.ListExchangesResponse, error) {
	if in.UserId == "" {
		return nil, gerrors.FailedPrecondition("missing-param.user_id", nil)
	}

	errParams := map[string]string{
		"user_id": in.UserId,
	}

	exchanges, err := dao.ListExchangesByUserID(ctx, in.UserId, in.GetActiveOnly())
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "exchanges_not_found_for_user_id"):
		return nil, gerrors.Augment(err, "failed_to_list_exchanges_by_user_id", errParams)
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_read_exchanges_by_user_id", errParams)
	}

	protos, err := marshaling.ExchangeDomainToProtos(exchanges)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_list_exchanges.at_least_one_exchange_has_an_unsupported_exchange_type", errParams)
	}

	return &accountproto.ListExchangesResponse{
		Exchanges: protos,
	}, nil
}
