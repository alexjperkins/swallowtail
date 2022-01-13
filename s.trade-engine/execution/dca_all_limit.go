package execution

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/monzo/slog"
	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/risk"
	or "swallowtail/s.trade-engine/orderrouter"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func init() {
	register(tradeengineproto.EXECUTION_STRATEGY_DCA_ALL_LIMIT, &DCAAllLimit{})
}

// DCAAllLimit executes a dca all limit strategy via direct market access.
type DCAAllLimit struct{}

// Execute ...
func (d *DCAAllLimit) Execute(
	ctx context.Context,
	strategy *tradeengineproto.TradeStrategy,
	participant *tradeengineproto.ExecuteTradeStrategyForParticipantRequest,
) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
	// Validation.
	switch {
	case len(strategy.Entries) < 2:
		return nil, gerrors.FailedPrecondition("failed_to_execute_dca_all_limit_strategy.invalid.more_than_one_entry_required", map[string]string{
			"trade_strategy_id": strategy.TradeStrategyId,
		})
	case participant.Venue == tradeengineproto.VENUE_UNREQUIRED:
		return nil, gerrors.FailedPrecondition("dca_all_limit.venue_required", nil)
	case participant.Risk == 0:
		return nil, gerrors.FailedPrecondition("dca_all_limit.participant_nil_risk", nil)
	}

	// Fetch venue specific credentials.
	venueCredentials, err := readVenueCredentials(ctx, participant.UserId, participant.Venue)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_dca_all_limit_strategy", nil)
	}

	// Read account balance.
	venueAccountBalance, err := readVenueAccountBalance(ctx, participant.Venue, strategy.InstrumentType, venueCredentials)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_dca_all_limit_strategy", nil)
	}

	// Calculate number of DCA positions.
	numberOfPositions := calculateNumberOfDCABuys(venueAccountBalance)

	// Marshal proto entries.
	var entries = make([]float64, 0, len(strategy.Entries))
	for _, e := range strategy.Entries {
		entries = append(entries, float64(e))
	}

	// Calculate positions.
	positions, err := risk.CalculatePositionsByRisk(entries, float64(strategy.StopLoss), numberOfPositions, strategy.TradeSide, tradeengineproto.DCA_EXECUTION_STRATEGY_LINEAR)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_dca_all_limit_strategy.calculate_risk", nil)
	}

	// Calculate total quantity/size from positions.
	totalQuantity := calculateTotalQuantityFromPositions(venueAccountBalance, float64(participant.Risk), positions)

	// Validate order against risk appetite constraints.
	if err := isTradeStrategyParticipantOverRiskAppetite(venueAccountBalance, totalQuantity); err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_dca_all_limit_strategy", map[string]string{
			"total_quantity": fmt.Sprintf("%f", totalQuantity),
			"venue_balance":  fmt.Sprintf("%f", venueAccountBalance),
		})
	}

	if err := isEnoughAvailableVenueMargin(venueAccountBalance); err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_dca_all_limit_strategy", map[string]string{
			"venue_balance":           fmt.Sprintf("%f", venueAccountBalance),
			"venue_min_margain_limit": fmt.Sprintf("%d", retailMinVenueMargainInUSDT),
		})
	}

	var (
		orders []*tradeengineproto.Order
		now    = time.Now().UTC()
	)

	errParams := map[string]string{
		"created_timetstamp":  now.String(),
		"number_of_positions": strconv.Itoa(numberOfPositions),
		"with_stop_loss":      strconv.FormatBool(strategy.StopLoss != 0),
		"risk":                fmt.Sprintf("%.02f", participant.Risk),
		"user_id":             participant.UserId,
		"asset":               strategy.Asset,
		"pair":                strategy.Pair.String(),
		"instrument":          strategy.Instrument,
		"venue":               participant.Venue.String(),
		"total_quantity":      fmt.Sprintf("%f", totalQuantity),
	}

	var exitTradeSide tradeengineproto.TRADE_SIDE
	switch strategy.TradeSide {
	case tradeengineproto.TRADE_SIDE_BUY, tradeengineproto.TRADE_SIDE_LONG:
		exitTradeSide = tradeengineproto.TRADE_SIDE_SELL
	default:
		exitTradeSide = tradeengineproto.TRADE_SIDE_BUY
	}

	// Add stop loss order.
	switch {
	case strategy.StopLoss == 0 && strategy.InstrumentType == tradeengineproto.INSTRUMENT_TYPE_FUTURE_PERPETUAL:
		slog.Warn(ctx, "Participant executing trade strategy without a stop loss: %s, %s", strategy.TradeStrategyId, participant.UserId)

		// Warn user of **not** using a stop loss. Best effort.
		if err := notifyUser(ctx, fmt.Sprintf("[%s] participant placing without a stop loss", strategy.ExecutionStrategy), participant.UserId); err != nil {
			slog.Error(ctx, "Failed to notifiy user: %v", err)
		}
	default:
		orders = append(orders, &tradeengineproto.Order{
			ActorId:          tradeengineproto.TradeEngineActorSatoshiSystem,
			Asset:            strategy.Asset,
			Pair:             strategy.Pair,
			Instrument:       strategy.Instrument,
			InstrumentType:   strategy.InstrumentType,
			OrderType:        tradeengineproto.ORDER_TYPE_STOP_MARKET,
			TradeSide:        exitTradeSide,
			StopPrice:        strategy.StopLoss,
			Quantity:         float32(totalQuantity),
			ReduceOnly:       true,
			WorkingType:      tradeengineproto.WORKING_TYPE_MARK_PRICE,
			Venue:            participant.Venue,
			CreatedTimestamp: now.Unix(),
		})
	}

	// Add entry orders.
	for _, p := range positions {
		orders = append(orders, &tradeengineproto.Order{
			ActorId:          tradeengineproto.TradeEngineActorSatoshiSystem,
			Instrument:       strategy.Instrument,
			Asset:            strategy.Asset,
			Pair:             strategy.Pair,
			InstrumentType:   strategy.InstrumentType,
			OrderType:        tradeengineproto.ORDER_TYPE_LIMIT,
			TradeSide:        strategy.TradeSide,
			LimitPrice:       float32(p.Price),
			Quantity:         float32(venueAccountBalance) * participant.Risk * float32(p.RiskCoefficient) / 100,
			WorkingType:      tradeengineproto.WORKING_TYPE_MARK_PRICE,
			Venue:            participant.Venue,
			CreatedTimestamp: now.Unix(),
			TimeInForce:      tradeengineproto.TIME_IN_FORCE_GOOD_TILL_CANCELLED,
		})
	}

	// Add take profits.
	tps := calculateTakeProfits(totalQuantity, strategy.TakeProfits)
	for _, tp := range tps {
		orders = append(orders, &tradeengineproto.Order{
			ActorId:          tradeengineproto.TradeEngineActorSatoshiSystem,
			Instrument:       strategy.Instrument,
			Asset:            strategy.Asset,
			Pair:             strategy.Pair,
			InstrumentType:   strategy.InstrumentType,
			OrderType:        tradeengineproto.ORDER_TYPE_TAKE_PROFIT_MARKET,
			TradeSide:        exitTradeSide,
			StopPrice:        float32(tp.StopPrice),
			Quantity:         float32(tp.Quantity),
			WorkingType:      tradeengineproto.WORKING_TYPE_MARK_PRICE,
			Venue:            participant.Venue,
			ReduceOnly:       true,
			CreatedTimestamp: now.Unix(),
		})
	}

	// Execute orders sequentially; gather successful orders, here we return early on the first failed order.
	// Here we manage risk, by placing the stop first - this is the most important.
	var successfulOrders = make([]*tradeengineproto.Order, 0, len(orders))
	for i, o := range orders {
		successfulOrder, err := or.RouteAndExecuteNewOrder(ctx, o, participant.Venue, strategy.InstrumentType, venueCredentials)
		if err != nil {
			slog.Error(ctx, "Failed to execute given order: %+v, Error: %v", o, err, errParams)
			return &tradeengineproto.ExecuteTradeStrategyForParticipantResponse{
				NotionalSize:           float32(totalQuantity),
				Venue:                  participant.Venue,
				NumberOfExecutedOrders: int64(i),
				ExecutionStrategy:      strategy.ExecutionStrategy,
				SuccessfulOrders:       successfulOrders,
				Timestamp:              timestamppb.Now(),
				Error: &tradeengineproto.ExecutionError{
					ErrorMessage: gerrors.Augment(err, "failed_to_execute_dca_all_limit_strategy", nil).Error(),
					FailedOrder:  o,
				},
			}, nil
		}

		slog.Info(ctx, "Order placed: %s [%s] %s", successfulOrder.Venue, successfulOrder.ExternalOrderId, successfulOrder.Instrument)
		successfulOrders = append(successfulOrders, successfulOrder)
	}

	slog.Info(ctx, "Successfully placed dca all limit trade strategy: %s for user: %s, risk: , total quantity: ", strategy.TradeStrategyId, participant.UserId, participant.Risk, totalQuantity)

	// TODO: Store into persistance layer.

	return &tradeengineproto.ExecuteTradeStrategyForParticipantResponse{
		NotionalSize:           float32(totalQuantity),
		NumberOfExecutedOrders: int64(len(successfulOrders)),
		ExecutionStrategy:      strategy.ExecutionStrategy,
		SuccessfulOrders:       successfulOrders,
		Timestamp:              timestamppb.Now(),
		Venue:                  participant.Venue,
		Asset:                  strategy.Asset,
		Pair:                   strategy.Pair,
		TradeParticipantId:     participant.UserId,
		Instrument:             strategy.Instrument,
	}, nil
}
