package execution

import (
	"context"
	"fmt"
	"sync"

	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

var (
	executionRegistry map[tradeengineproto.EXECUTION_STRATEGY]StrategyExecution
	executionMu       sync.RWMutex
)

// StrategyExecution defines the execution execution.
type StrategyExecution interface {
	Execute(ctx context.Context, strategy *tradeengineproto.TradeStrategy, participant *tradeengineproto.ExecuteTradeStrategyForParticipantRequest) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error)
}

// ExecuteTradeStrategy executes the given trade strategy with the given execution executionrithm.
func ExecuteTradeStrategy(ctx context.Context, strategy *tradeengineproto.TradeStrategy, participant *tradeengineproto.ExecuteTradeStrategyForParticipantRequest) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
	errParams := map[string]string{
		"execution_strategy": strategy.ExecutionStrategy.String(),
		"user_id":            participant.UserId,
		"actor_id":           participant.ActorId,
		"venue":              participant.Venue.String(),
	}

	execution, ok := getStrategyExecution(strategy.ExecutionStrategy)
	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_execute_trading_strategy.invalid_execution_strategy", errParams)
	}

	rsp, err := execution.Execute(ctx, strategy, participant)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_trading_strategy_execution_for_participant", errParams)
	}

	return rsp, nil
}

func register(strategy tradeengineproto.EXECUTION_STRATEGY, handler StrategyExecution) {
	executionMu.Lock()
	defer executionMu.Unlock()

	if _, ok := executionRegistry[strategy]; ok {
		panic(fmt.Sprintf("Failed to register execution strategy: strategy already registered; %s", strategy))
	}

	executionRegistry[strategy] = handler
}

func getStrategyExecution(strategy tradeengineproto.EXECUTION_STRATEGY) (StrategyExecution, bool) {
	executionMu.Lock()
	defer executionMu.Unlock()

	a, ok := executionRegistry[strategy]
	if !ok {
		return nil, false
	}

	return a, true
}
