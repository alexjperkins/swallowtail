package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/dao"
	"swallowtail/s.trade-engine/exchange"
	"swallowtail/s.trade-engine/marshaling"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// AddParticipantToTrade ...
func (s *TradeEngineService) AddParticipantToTrade(
	ctx context.Context, in *tradeengineproto.AddParticipantToTradeRequest,
) (*tradeengineproto.AddParticipantToTradeResponse, error) {
	switch {
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	case !isActorValid(in.ActorId):
		return nil, gerrors.Unauthenticated("failed_to_add_participant_to_trade.unauthorized", nil)
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.TradeId == "":
		return nil, gerrors.BadParam("missing_param.trade_id", nil)
	}

	errParams := map[string]string{
		"actor_id": in.ActorId,
		"trade_id": in.TradeId,
		"exchange": in.Exchange,
	}

	// Read trade to see if it exists.
	trade, err := dao.ReadTradeByTradeID(ctx, in.TradeId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade", errParams)
	}

	// Read trade participant to see if that already exists
	existingTradeParticipant, err := dao.ReadTradeParticipantByTradeID(ctx, trade.TradeID, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "not_found.trade_participant"):
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.failed_check_if_trade_participant_already_exists", errParams)
	}

	if existingTradeParticipant != nil {
		return nil, gerrors.AlreadyExists("failed_to_add_participant_to_trade.trade_already_exists", errParams)
	}

	// Validate our trade participant.
	if err := validateTradeParticipant(in, trade); err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.invalid_trade_participant", errParams)
	}

	// Read the users primary exchange; here we do an implicit account check. This ensures the user does have an account with us.
	// **we** don't need to verify if the user is a futures paying user.
	primaryExchangeCredentials, err := readPrimaryExchangeCredentials(ctx, in.UserId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.read_primary_exchange", errParams)
	}

	// Execute the trade.
	exchangeTradeRsp, err := exchange.ExecuteFuturesTradeForParticipant(ctx, trade, in, primaryExchangeCredentials)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.execute_trade", errParams)
	}

	// Embelish the trade participant before marshaling & persisting.
	in.Exchange = primaryExchangeCredentials.ExchangeType.String()
	in.Size = float32(exchangeTradeRsp.NotionalSize)

	tradeParticipant := marshaling.TradeParticipantProtoToDomain(in, exchangeTradeRsp.ExchangeTradeID, exchangeTradeRsp.ExecutionTimestamp)
	if err := dao.AddTradeParticpantToTrade(ctx, tradeParticipant); err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.dao", errParams)
	}

	// Read trade participant back out to get the trade participant id.
	existingTradeParticipant, err = dao.ReadTradeParticipantByTradeID(ctx, trade.TradeID, in.UserId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.failed_to_read_created_trade_participant", errParams)
	}

	return &tradeengineproto.AddParticipantToTradeResponse{
		ExchangeTradeId:    exchangeTradeRsp.ExchangeTradeID,
		TradeParticipantId: existingTradeParticipant.TradeParticipantID,
		TradeId:            trade.TradeID,
		NotionalSize:       float32(exchangeTradeRsp.NotionalSize),
		Timestamp:          timestamppb.New(exchangeTradeRsp.ExecutionTimestamp),
		Asset:              trade.Asset,
		Exchange:           tradeParticipant.Exchange,
	}, nil
}
