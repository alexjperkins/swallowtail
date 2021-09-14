package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func executeTradeForUser(ctx context.Context, userID, tradeID string, riskPercentage int) (*tradeengineproto.AddParticipantToTradeResponse, error) {
	rsp, err := (&tradeengineproto.AddParticipantToTradeRequest{
		ActorId: tradeengineproto.TradeEngineActorSatoshiSystem,
		IsBot:   true,
		UserId:  userID,
		TradeId: tradeID,
		Risk:    float32(riskPercentage),
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_trade_for_user", nil)
	}

	return rsp, nil
}
