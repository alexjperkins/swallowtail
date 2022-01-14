package execution

import (
	"context"

	"swallowtail/libraries/gerrors"
	or "swallowtail/s.trade-engine/orderrouter"
	tradeengineproto "swallowtail/s.trade-engine/proto"

	"github.com/monzo/slog"
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

func executeOrdersSequentiallyWithoutRetry(
	ctx context.Context,
	orders []*tradeengineproto.Order,
	venue tradeengineproto.VENUE,
	instrumentType tradeengineproto.INSTRUMENT_TYPE,
	credentials *tradeengineproto.VenueCredentials,
) ([]*tradeengineproto.Order, *tradeengineproto.ExecutionError) {
	var successfulOrders = make([]*tradeengineproto.Order, 0, len(orders))
	for _, order := range orders {
		successfulOrder, err := or.RouteAndExecuteNewOrder(ctx, order, venue, instrumentType, credentials)
		if err != nil {
			slog.Error(ctx, "Failed to execute given order: %+v, Error: %v", order, err)
			return successfulOrders, &tradeengineproto.ExecutionError{
				ErrorMessage: gerrors.Augment(err, "failed_to_execute_order", nil).Error(),
				FailedOrder:  order,
			}
		}

		slog.Info(ctx, "Order placed: %s [%s] %s", successfulOrder.Venue, successfulOrder.ExternalOrderId, successfulOrder.Instrument)
		successfulOrders = append(successfulOrders, successfulOrder)
	}

	return successfulOrders, nil
}
