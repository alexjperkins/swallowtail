package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func executeTradeStrategyForParticipant(
	ctx context.Context,
	userID,
	tradeStrategyID string,
	riskPercentage int,
) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
	// Fetch primary exchange for the user to use.
	// TODO: this probably should be moved upstream.
	venue, err := getPrimaryVenueByUserID(ctx, userID)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_trade_strategy_for_user", nil)
	}

	rsp, err := (&tradeengineproto.ExecuteTradeStrategyForParticipantRequest{
		Venue:           venue,
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

func getPrimaryVenueByUserID(ctx context.Context, userID string) (tradeengineproto.VENUE, error) {
	rsp, err := (&accountproto.ReadPrimaryVenueAccountByUserIDRequest{
		UserId:  userID,
		ActorId: accountproto.ActorSystemTradeEngine,
	}).Send(ctx).Response()
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_read_primary_venue_by_user_id", nil)
	}

	return rsp.PrimaryVenueAccount.Venue, nil
}
