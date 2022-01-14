package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	ftxproto "swallowtail/s.ftx/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func isValidActor(ctx context.Context, actorID string) (bool, error) {
	if actorID == ftxproto.FTXDepositAccountActorPaymentsSystem {
		return true, nil
	}

	// Check the actor is authorized to make this request.
	account, err := (&accountproto.ReadAccountRequest{
		UserId: actorID,
	}).Send(ctx).Response()
	if err != nil {
		return false, gerrors.Augment(err, "failed_to_list_account_deposits.failed_to_read_account_of_actor", nil)
	}

	return account.GetAccount().IsAdmin, nil
}

func validateCredentials(credentials *tradeengineproto.VenueCredentials) error {
	switch {
	case credentials == nil:
		return gerrors.BadParam("missing_param.credentials", nil)
	case credentials.ApiKey == "":
		return gerrors.BadParam("missing_param.credentials.api_key", nil)
	case credentials.SecretKey == "":
		return gerrors.BadParam("missing_param.credentials.secret_key", nil)
	case credentials.Subaccount == "":
		return gerrors.BadParam("missing_param.credentials.subaccount", nil)
	default:
		return nil
	}
}

func validateOrder(order *tradeengineproto.Order) error {
	switch {
	case order.Venue != tradeengineproto.VENUE_FTX:
		return gerrors.FailedPrecondition("invalid_venue.expecting_ftx", nil)
	case order.Instrument == "" && order.Asset == "":
		return gerrors.BadParam("missing_param.instrument_or_asset", nil)
	case order.ClosePosition && order.Quantity == 0:
		return gerrors.Unimplemented("unimplemented.close_position", nil)
	case order.Quantity <= 0:
		return gerrors.BadParam("missing_param.quantity", nil)
	case order.InstrumentType == tradeengineproto.INSTRUMENT_TYPE_FORWARD:
		return gerrors.Unimplemented("instrument_type.forward", nil)
	}

	switch order.OrderType {
	case tradeengineproto.ORDER_TYPE_LIMIT:
		switch {
		case order.LimitPrice <= 0:
			return gerrors.BadParam("bad_param.limit_price", nil)
		}
	case tradeengineproto.ORDER_TYPE_MARKET, tradeengineproto.ORDER_TYPE_TAKE_PROFIT_MARKET:
		switch {
		case order.StopPrice <= 0:
			return gerrors.BadParam("bad_param.stop_price", nil)
		}
	case tradeengineproto.ORDER_TYPE_STOP_LIMIT, tradeengineproto.ORDER_TYPE_TAKE_PROFIT_LIMIT:
		switch {
		case order.LimitPrice <= 0:
			return gerrors.BadParam("bad_param.limit_price", nil)
		case order.StopPrice <= 0:
			return gerrors.BadParam("bad_param.stop_price", nil)
		}
	}

	return nil
}
