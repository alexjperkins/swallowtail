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
	default:
		return nil
	}
}

func validateOrder(order *tradeengineproto.Order) error {
	switch {
	case order.Instrument== "":
		return gerrors.BadParam("bad_param.symbol", nil)
	}

	case order.Price < 0:
		return gerrors.BadParam("bad_param.price.negative", nil)
	case order.TriggerPrice < 0:
		return gerrors.BadParam("bad_param.trigger_price.negative", nil)
	case order.OrderPrice < 0:
		return gerrors.BadParam("bad_param.order_price.negative", nil)
	case order.TrailValue < 0:
		return gerrors.BadParam("bad_param.trail_value.negative", nil)
	case order.Quantity == 0:
		return gerrors.BadParam("missing_param.quantity", nil)
	}

	switch {
	case order.Type != ftxproto.FTX_TRADE_TYPE_FTX_TRADE_TYPE_MARKET && order.Price == 0:
		return gerrors.BadParam("missing_param.price", map[string]string{
			"type": order.Type.String(),
		})
	}

	switch order.Type {
	case ftxproto.FTX_TRADE_TYPE_FTX_TRADE_TYPE_STOP,
		ftxproto.FTX_TRADE_TYPE_FTX_TRADE_TYPE_TAKE_PROFIT:
		if order.TriggerPrice == 0 && order.OrderPrice == 0 {
			return gerrors.BadParam("missing_param.trigger_price_or_order_price", map[string]string{
				"type": order.Type.String(),
			})
		}

	case ftxproto.FTX_TRADE_TYPE_FTX_TRADE_TYPE_TRALING_STOP:
		if order.TrailValue == 0 {
			return gerrors.BadParam("missing_param.trail_value", map[string]string{
				"type": order.Type.String(),
			})
		}
	}

	return nil
}
