package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"

	"github.com/monzo/terrors"
)

func (s *AccountService) ListExchanges(
	ctx context.Context, in *accountproto.ListExchangesRequest,
) (*accountproto.ListExchangesResponse, error) {
	if in.UserId == "" {
		return nil, terrors.PreconditionFailed("missing-param.user_id", "Cannot list exchanges; missing user id", nil)
	}

	errParams := map[string]string{
		"user_id": in.UserId,
	}

	exchanges, err := dao.ListExchangesByUserID(ctx, in.UserId, in.GetActiveOnly())
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "exchanges_not_found_for_user_id"):
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_read_exchanges_by_user_id", errParams)
	}

	protos, err := marshaling.ExchangeDomainToProtos(exchanges)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to list exchanges; at least one  exchange has an unsupported exchange type", errParams)
	}

	return &accountproto.ListExchangesResponse{
		Exchanges: protos,
	}, nil
}
