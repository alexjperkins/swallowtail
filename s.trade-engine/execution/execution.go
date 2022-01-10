package execution

import (
	"context"

	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// StrategyExecution defines the execution execution.
type StrategyExecution interface {
	Execute(
		ctx context.Context,
		strategy *tradeengineproto.TradeStrategy,
		participant *tradeengineproto.ExecuteTradeStrategyForParticipantRequest,
	) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error)
}

// ExecuteTradeStrategyForParticipant executes the given trade strategy with the given execution algorithm.
func ExecuteTradeStrategyForParticipant(
	ctx context.Context,
	strategy *tradeengineproto.TradeStrategy,
	participant *tradeengineproto.ExecuteTradeStrategyForParticipantRequest,
) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
	errParams := map[string]string{
		"execution_strategy": strategy.ExecutionStrategy.String(),
		"user_id":            participant.UserId,
		"actor_id":           participant.ActorId,
		"venue":              participant.Venue.String(),
	}

	// Fetch strategy from local registry.
	executionStrategy, ok := getStrategyExecution(strategy.ExecutionStrategy)
	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_execute_trading_strategy.invalid_execution_strategy", errParams)
	}

	// Execute.
	rsp, err := executionStrategy.Execute(ctx, strategy, participant)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_trading_strategy_execution_for_participant", errParams)
	}

	return rsp, nil
}
