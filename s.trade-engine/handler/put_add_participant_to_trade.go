package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/dao"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// AddParticipantToTrade ...
func (s *TradeEngineService) AddParticipantToTrade(
	ctx context.Context, in *tradeengineproto.AddParticipantToTradeRequest,
) (*tradeengineproto.AddParticipantToTradeRequest, error) {
	switch {
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	case !isActorValid(in.ActorId):
		return nil, gerrors.Unauthenticated("failed_to_add_participant_to_trade.unauthorized", nil)
	case in.TradeId == "":
		return nil, gerrors.BadParam("missing_param.trade_id", nil)
	case in.Exchange == "":
		return nil, gerrors.BadParam("missing_param.exchange", nil)
	case in.Risk <= 0:
		return nil, gerrors.BadParam("bad_param.risk_cannot_be_zero_or_below", nil)
	}

	errParams := map[string]string{
		"actor_id": in.ActorId,
		"trade_id": in.TradeId,
		"exchange": in.Exchange,
	}

	// Read trade to see if it exists
	trade, err := dao.ReadTradeByTradeID(ctx, in.TradeId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade", errParams)
	}

	if err := validateTradeParticipant(in, trade); err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.invalid_trade_participant", errParams)
	}

	// Read the users primary exchange; here we do an implicit account check. This ensures the user does have an account with us.
	// **we** don't need to verify if the user is a futures paying user.
	primaryExchangeCredentials, err := readPrimaryExchangeCredentials(ctx, in.UserId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.read_primary_exchange", errParams)
	}

	// TODO
	// parse credentials into a request & execute.
	// on success receipt we can store in our database with the exchange trade id.

	return &tradeengineproto.AddParticipantToTradeRequest{}, nil
}
