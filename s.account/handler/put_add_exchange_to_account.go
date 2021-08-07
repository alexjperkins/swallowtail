package handler

import (
	"context"
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
	errParams := map[string]string{
		"user_id": in.Exchange.UserId,
	}

	exchange, err := marshaling.ExchangeProtoToDomain(in.Exchange)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to marshal request", errParams)
	}

	if err := dao.AddExchange(ctx, exchange); err != nil {
		return nil, terrors.Augment(err, "Failed to add exchange to account.", errParams)
	}

	slog.Info(ctx, "Added new exchange to account", errParams)

	return &accountproto.AddExchangeResponse{}, nil
}
