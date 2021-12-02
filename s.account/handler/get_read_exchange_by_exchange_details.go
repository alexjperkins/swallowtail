package handler

import (
	"context"
	"strings"
	"swallowtail/libraries/gerrors"

	"swallowtail/s.account/dao"
	"swallowtail/s.account/marshaling"
	accountproto "swallowtail/s.account/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// ReadExchangeByExchangeDetails ...
func (s *AccountService) ReadExchangeByExchangeDetails(
	ctx context.Context, in *accountproto.ReadExchangeByExchangeDetailsRequest,
) (*accountproto.ReadExchangeByExchangeDetailsResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	case in.Exchange == "":
		return nil, gerrors.BadParam("missing_param.exchange", nil)
	case in.RequestContext == "":
		return nil, gerrors.BadParam("missing_param.request_context", nil)
	}

	errParams := map[string]string{
		"user_id":    in.UserId,
		"actor_id":   in.ActorId,
		"exchange":   in.Exchange,
		"subaccount": in.Subaccount,
	}

	// Validate the request for exchange credentials.
	if err := validateExchangeDetailsRequest(in); err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_exchange_by_exchange_details", errParams)
	}

	// Read.
	exchange, err := dao.ReadExchangeByExchangeDetails(ctx, in.Exchange, in.UserId, in.Subaccount)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_exchange_by_exchange_details.dao", errParams)
	}

	// Marshal.
	exchangeProto, err := marshaling.ExchangeDomainToProtoUnmasked(exchange)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_exchange_by_exchange_details.marshal_to_proto", errParams)
	}

	return &accountproto.ReadExchangeByExchangeDetailsResponse{
		Exchange: exchangeProto,
	}, nil
}

func validateExchangeDetailsRequest(in *accountproto.ReadExchangeByExchangeDetailsRequest) error {
	switch in.RequestContext {
	case accountproto.RequestContextOrderRequest, accountproto.RequestContextUserRequest:
	default:
		return gerrors.FailedPrecondition("invalid_request_context", nil)
	}

	switch strings.ToUpper(in.Exchange) {
	case tradeengineproto.VENUE_BINANCE.String(), tradeengineproto.VENUE_FTX.String(), tradeengineproto.VENUE_BITFINEX.String(), tradeengineproto.VENUE_DERIBIT.String():
	default:
		return gerrors.Unimplemented("exchange.unimplemented", nil)
	}

	switch in.ActorId {
	case in.UserId:
	case accountproto.ActorSystemTradeEngine, accountproto.ActorSystemPayments:
	default:
		return gerrors.FailedPrecondition("bad_actor", nil)
	}

	return nil
}
