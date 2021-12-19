package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func executeTradeStrategyForParticipant(ctx context.Context, userID, tradeStrategyID string, riskPercentage int) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
	rsp, err := (&tradeengineproto.ExecuteTradeStrategyForParticipantRequest{
		ActorId:         tradeengineproto.TradeEngineActorSatoshiSystem,
		IsBot:           true,
		UserId:          userID,
		TradeStrategyId: tradeStrategyID,
		Risk:            float32(riskPercentage),
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_trade_strategy_for_user", nil)
	}

	return rsp, nil
}
