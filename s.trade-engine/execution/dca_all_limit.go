package execution

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/risk"
	or "swallowtail/s.trade-engine/orderrouterv2"
	tradeengineproto "swallowtail/s.trade-engine/proto"

	"github.com/monzo/slog"
)

func init() {
	register(tradeengineproto.EXECUTION_STRATEGY_DCA_ALL_LIMIT, &DCAAllLimit{})
}

type DCAAllLimit struct{}

func (d *DCAAllLimit) Execute(ctx context.Context, strategy *tradeengineproto.TradeStrategy, participant *tradeengineproto.ExecuteTradeStrategyForParticipantRequest) (*tradeengineproto.ExecuteTradeStrategyForParticipantResponse, error) {
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
	positions, err := risk.CalculatePositionsByRisk(entries, float64(strategy.StopLoss), float64(participant.Risk), numberOfPositions, strategy.TradeSide, tradeengineproto.DCA_EXECUTION_STRATEGY_LINEAR)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_dca_all_limit_strategy.calculate_risk", nil)
	}

	// Calculate total quantity/size from positions.
	totalQuantity := calculateTotalQuantityFromPositions(venueAccountBalance, float64(participant.Risk), positions)

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
		"venue":               participant.Venue.String(),
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
		if err := notifyUser(ctx, "DCA_ALL_LIMIT: Placing without a STOP LOSS", participant.UserId); err != nil {
			slog.Error(ctx, "Failed to notifiy user: %v", err)
		}
	default:

		orders = append(orders, &tradeengineproto.Order{
			ActorId:          tradeengineproto.TradeEngineActorSatoshiSystem,
			Instrument:       fmt.Sprintf("%s%s", strategy.Asset, strategy.Pair.String()),
			Pair:             strategy.Pair.String(),
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
			Instrument:       fmt.Sprintf("%s%s", strategy.Asset, strategy.Pair.String()),
			Pair:             strategy.Pair.String(),
			InstrumentType:   strategy.InstrumentType,
			OrderType:        tradeengineproto.ORDER_TYPE_LIMIT,
			TradeSide:        strategy.TradeSide,
			LimitPrice:       float32(p.Price),
			Quantity:         float32(venueAccountBalance) * participant.Risk * float32(p.RiskCoefficient),
			WorkingType:      tradeengineproto.WORKING_TYPE_MARK_PRICE,
			Venue:            participant.Venue,
			CreatedTimestamp: now.Unix(),
		})
	}

	tps := calculateTakeProfits(totalQuantity, strategy.TakeProfits)

	// Add take profits.
	for _, tp := range tps {
		orders = append(orders, &tradeengineproto.Order{
			ActorId:          tradeengineproto.TradeEngineActorSatoshiSystem,
			Instrument:       fmt.Sprintf("%s%s", strategy.Asset, strategy.Pair.String()),
			Pair:             strategy.Pair.String(),
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

	successfulOrders, err := or.RouteExecuteNewOrder(ctx, orders, participant.Venue, strategy.InstrumentType, venueCredentials)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_dca_all_limit_strategy.orderrouter", errParams)
	}

	// Store DB.

	return &tradeengineproto.ExecuteTradeStrategyForParticipantResponse{
		SuccessfulOrders: successfulOrders,
		Exchange:         participant.Venue.String(), // TODO: change
	}, nil
}
