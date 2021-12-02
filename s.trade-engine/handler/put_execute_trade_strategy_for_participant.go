package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	"swallowtail/s.trade-engine/dao"
	"swallowtail/s.trade-engine/marshaling"
	"swallowtail/s.trade-engine/orderrouter"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// ExecuteTradeStrategyForParticipant ...
func (s *TradeEngineService) ExecuteTradeStrategyForParticipant(
	ctx context.Context, in *tradeengineproto.ExecuteTradeStrategyForParticipantRequest,
) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
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
		"venue":    in.Venue.String(),
	}

	// Read trade strategy to see if it exists.
	trade, err := dao.ReadTradeStrategyByTradeStrategyID(ctx, in.TradeId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade", errParams)
	}

	// Read trade participant to see if that already exists
	existingTradeParticipant, err := dao.ReadTradeStrategyParticipantByTradeStrategyID(ctx, trade.TradeStrategyID, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "not_found.trade_participant"):
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.failed_check_if_trade_participant_already_exists", errParams)
	}

	if existingTradeParticipant != nil {
		return nil, gerrors.AlreadyExists("failed_to_add_participant_to_trade.trade_already_exists", errParams)
	}

	// Validate our trade strategy participant.
	if err := validateTradeStrategyParticipant(in, trade); err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.invalid_trade_participant", errParams)
	}

	// Read the users primary exchange; here we do an implicit account check. This ensures the user does have an account with us.
	// **we** don't need to verify if the user is a futures paying user.
	primaryExchangeCredentials, err := readPrimaryExchangeCredentials(ctx, in.UserId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.read_primary_exchange", errParams)
	}

	// Execute the trade.
	exchangeTradeRsp, err := orderrouter.ExecuteFuturesTradeStrategyForParticipant(ctx, trade, in, primaryExchangeCredentials)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.execute_trade", errParams)
	}

	// Translate the exchange.
	var venue tradeengineproto.VENUE
	switch primaryExchangeCredentials.ExchangeType {
	case accountproto.ExchangeType_BINANCE:
	case accountproto.ExchangeType_BITFINEX:
	case accountproto.ExchangeType_FTX:
	case accountproto.ExchangeType_DERIBIT:
	default:
		return nil, gerrors.Unimplemented("failed_to_add_participant_to_trade_strategy.unimplemented_exchange.credentials", map[string]string{
			"exchange": primaryExchangeCredentials.ExchangeType.String(),
		})
	}

	// Embelish the trade participant before marshaling & persisting.
	in.Venue = venue
	in.Size = float32(exchangeTradeRsp.NotionalSize)

	tradeStrategyParticipant := marshaling.TradeParticipantProtoToDomain(in, exchangeTradeRsp.ExchangeTradeIDs, exchangeTradeRsp.ExecutionTimestamp)
	if err := dao.AddParticpantToTradeStrategy(ctx, tradeStrategyParticipant); err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.dao", errParams)
	}

	// Read trade participant back out to get the trade participant id.
	existingTradeParticipant, err = dao.ReadTradeStrategyParticipantByTradeStrategyID(ctx, trade.TradeStrategyID, in.UserId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.failed_to_read_created_trade_participant", errParams)
	}

	return &tradeengineproto.ExecuteTradeStrategyForParticipantResponse{
		ExchangeTradeIds:       exchangeTradeRsp.ExchangeTradeIDs,
		TradeParticipantId:     existingTradeParticipant.TradeParticipantID,
		TradeId:                trade.TradeID,
		NotionalSize:           float32(exchangeTradeRsp.NotionalSize),
		Timestamp:              timestamppb.New(exchangeTradeRsp.ExecutionTimestamp),
		Asset:                  trade.Asset,
		Exchange:               tradeParticipant.Exchange,
		NumberOfExecutedOrders: int64(exchangeTradeRsp.NumberOfExecutedOrders),
		ExecutionAlgoStrategy:  exchangeTradeRsp.ExecutionAlgoStrategy,
	}, nil
}
