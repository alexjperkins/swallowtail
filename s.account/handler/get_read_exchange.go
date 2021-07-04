package handler

import (
	"context"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"

	"github.com/monzo/terrors"
)

// ReadExchange ...
func (s *AccountService) ReadExchange(
	ctx context.Context, in *accountproto.ReadExchangeRequest,
) (*accountproto.ReadExchangeResponse, error) {
	if in.ExchangeId == "" {
		return nil, terrors.PreconditionFailed("missing-param.exchange-id", "Cannot read exchange; missing exchange id", nil)

	}

	errParams := map[string]string{
		"exchange_id": in.ExchangeId,
	}

	exchange, err := dao.ReadExchangeByExchangeID(ctx, in.ExchangeId)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to read exchange", errParams)
	}

	proto, err := marshaling.ExchangeDomainToProto(exchange)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to read exchange", errParams)
	}

	return &accountproto.ReadExchangeResponse{
		Exchange: proto,
	}, nil
}
