package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/dao"
	"swallowtail/s.trade-engine/execution"
	"swallowtail/s.trade-engine/marshaling"
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
	tradeStrategy, err := dao.ReadTradeStrategyByTradeStrategyID(ctx, in.TradeId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade", errParams)
	}

	// Read trade participant to see if that already exists.
	existingTradeParticipant, err := dao.ReadTradeStrategyParticipantByTradeStrategyID(ctx, tradeStrategy.TradeStrategyID, in.UserId)
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound, "not_found.trade_participant"):
	case err != nil:
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.failed_check_if_trade_participant_already_exists", errParams)
	case existingTradeParticipant != nil:
		return nil, gerrors.AlreadyExists("failed_to_add_participant_to_trade.trade_already_exists", errParams)
	}

	// Validate our trade strategy participant.
	if err := validateTradeStrategyParticipant(in, tradeStrategy); err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.invalid_trade_participant", errParams)
	}

	// Marshal domain trade strategy to proto; here we can leverage enums over order parameters.
	tradeStrategyProto := marshaling.TradeStrategyDomainStrategyToProto(tradeStrategy)

	// Execute the trade.
	rsp, err := execution.ExecuteTradeStrategyForParticipant(ctx, tradeStrategyProto, in)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_add_participant_to_trade.execute_trade", errParams)
	}

	return rsp, nil
}
