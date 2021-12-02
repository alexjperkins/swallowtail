package execution

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/risk"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func init() {
	register(tradeengineproto.EXECUTION_STRATEGY_DCA_ALL_LIMIT, &DCAAllLimit{})
}

type DCAAllLimit struct{}

func (d *DCAAllLimit) Execute(ctx context.Context, strategy *tradeengineproto.TradeStrategy, participant *tradeengineproto.ExecuteTradeStrategyForParticipantRequest) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
	venueCredentials, err := readVenueCredentials(ctx, participant.UserId, participant.Venue)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_dca_all_limit_strategy", nil)
	}

	venueAccountBalance, err := readVenueAccountBalance(ctx, participant.Venue, venueCredentials)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_dca_all_limit_strategy", nil)
	}

	// TODO: 7 move to settings or something for a given user or even a constant.
	positions, err := risk.CalculatePositionsByRisk(strategy.Entries, strategy.StopLoss, participant.Risk, 7, strategy.TradeSide, tradeengineproto.DCA_EXECUTION_STRATEGY_LINEAR)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_dca_all_limit_strategy.calculate_risk", nil)
	}

	return nil, nil
}
