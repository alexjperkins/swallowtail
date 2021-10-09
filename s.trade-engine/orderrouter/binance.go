package orderrouter

import (
	"context"
	"fmt"
	"strings"

	"github.com/monzo/slog"
	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/libraries/gerrors"
	riskutil "swallowtail/libraries/risk"
	accountproto "swallowtail/s.account/proto"
	binanceproto "swallowtail/s.binance/proto"
	"swallowtail/s.trade-engine/domain"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

const (
	defaultBinanceDCAOrders = 5
)

func calculuateNotionalSizesFromBinanceAccount(entries []float64, stopLoss, risk, accountSize float64, side *tradeengineproto.TRADE_SIDE, strategy *tradeengineproto.DCA_STRATEGY) ([]float64, error) {
	risks, err := riskutil.CalculateNotionalSizesFromPositionAndRisk(entries, stopLoss, risk, defaultBinanceDCAOrders, side, strategy)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_calculate_risk_sizes.failed_to_calculate_notional_size_from_risk", nil)
	}

	notionalSizes := make([]float64, 0, len(entries))
	for _, risk := range risks {
		notionalSizes = append(notionalSizes, accountSize*risk)
	}

	return notionalSizes, nil
}

func readBinancePerpetualFuturesAccountSize(ctx context.Context, credentials *binanceproto.Credentials) (float64, error) {
	rsp, err := (&binanceproto.ReadPerpetualFuturesAccountRequest{
		ActorId:     binanceproto.BinanceAccountActorTradeEngineSystem,
		Credentials: credentials,
	}).Send(ctx).Response()
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_read_binance_perpetual_futures_account", nil)
	}

	return float64(rsp.AvailableBalance), nil
}

func readAccountByUserID(ctx context.Context, userID string) (*accountproto.Account, error) {
	rsp, err := (&accountproto.ReadAccountRequest{
		UserId: userID,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_account_by_user_id", map[string]string{
			"user_id": userID,
		})
	}

	return rsp.Account, nil
}

func executeBinanceFuturesTrade(
	ctx context.Context,
	trade *domain.Trade,
	participant *tradeengineproto.AddParticipantToTradeRequest,
	credentials *binanceproto.Credentials,
) (*FuturesTradeResponse, error) {

	var notionalSizes []float64
	switch {
	case participant.Risk != 0:
		// Marshal side back into proto side.
		var side *tradeengineproto.TRADE_SIDE
		switch strings.ToLower(trade.TradeSide) {
		case "long":
			side = tradeengineproto.TRADE_SIDE_LONG.Enum()
		case "buy":
			side = tradeengineproto.TRADE_SIDE_BUY.Enum()
		case "short":
			side = tradeengineproto.TRADE_SIDE_SHORT.Enum()
		case "sell":
			side = tradeengineproto.TRADE_SIDE_SELL.Enum()
		default:
			return nil, gerrors.Unimplemented("failed_to_execute_binance_futures_trade.unimplemented_trade_side", map[string]string{
				"trade_side": trade.TradeSide,
			})
		}

		// Read binance perpetual futures account.
		binanceAccountSize, err := readBinancePerpetualFuturesAccountSize(ctx, credentials)
		if err != nil {
			return nil, err
		}

		// Read account.
		account, err := readAccountByUserID(ctx, participant.UserId)
		if err != nil {
			return nil, err
		}

		// Marshal default dca strategy.
		var dcaStrategy *tradeengineproto.DCA_STRATEGY
		switch account.DefaultDcaStrategy {
		case tradeengineproto.DCA_STRATEGY_CONSTANT.String():
			dcaStrategy = tradeengineproto.DCA_STRATEGY_CONSTANT.Enum()
		case tradeengineproto.DCA_STRATEGY_LINEAR.String():
			dcaStrategy = tradeengineproto.DCA_STRATEGY_LINEAR.Enum()
		case tradeengineproto.DCA_STRATEGY_EXPONENTIAL.String():
			dcaStrategy = tradeengineproto.DCA_STRATEGY_EXPONENTIAL.Enum()
		}

		// Calculate notional size of all orders.
		ns, err := calculuateNotionalSizesFromBinanceAccount(trade.Entries, trade.StopLoss, binanceAccountSize, float64(participant.Risk), side, dcaStrategy)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_execute_binance_futures_trade", map[string]string{
				"side": side.String(),
			})
		}
		notionalSizes = ns
	default:
		return nil, gerrors.Unimplemented("failed_to_execute_binance_futures_trade.notional_size_calc_unimplimented", nil)
	}

	totalNotionalSize := sum(notionalSizes)
	errParams := map[string]string{
		"total_notional_size": fmt.Sprintf("%v", totalNotionalSize),
		"total_risk":          fmt.Sprintf("%v", participant.Risk),
	}
	orders := make([]*binanceproto.PerpetualFuturesOrder, 0, len(trade.Entries)+1)

	// Add stop loss order. We add this first for safety.
	switch trade.StopLoss {
	case 0:
		slog.Warn(ctx, "Creating trade without a stop loss: [%v] %s", trade.ActorType, trade.ActorID)
	default:
		orders = append(orders, &binanceproto.PerpetualFuturesOrder{
			StopPrice:     float32(trade.StopLoss),
			OrderType:     binanceproto.BinanceOrderType_STOP_MARKET,
			ClosePosition: true,
		})
	}

	// Add entry positions as order.
	for i, entry := range trade.Entries {
		// Parse order type.
		var orderType binanceproto.BinanceOrderType
		switch {
		case trade.OrderType == tradeengineproto.ORDER_TYPE_MARKET.String():
			orderType = binanceproto.BinanceOrderType_MARKET
		case trade.OrderType == tradeengineproto.ORDER_TYPE_LIMIT.String():
			orderType = binanceproto.BinanceOrderType_LIMIT
		case trade.OrderType == tradeengineproto.ORDER_TYPE_DCA_ALL_LIMIT.String():
			orderType = binanceproto.BinanceOrderType_LIMIT
		case trade.OrderType == tradeengineproto.ORDER_TYPE_DCA_FIRST_MARKET_REST_LIMIT.String() && i == len(trade.Entries)-1:
			orderType = binanceproto.BinanceOrderType_MARKET
		default:
			slog.Warn(ctx, "Binance order router recieved trade with unrecognised order type: %s", trade.OrderType)
			orderType = binanceproto.BinanceOrderType_LIMIT
		}

		// Parse trade type.
		var side binanceproto.BinanceTradeSide
		switch trade.TradeSide {
		case tradeengineproto.TRADE_SIDE_BUY.String(), tradeengineproto.TRADE_SIDE_LONG.String():
			side = binanceproto.BinanceTradeSide_BUY
		case tradeengineproto.TRADE_SIDE_SELL.String(), tradeengineproto.TRADE_SIDE_SHORT.String():
			side = binanceproto.BinanceTradeSide_SELL
		default:
			return nil, gerrors.FailedPrecondition("failed_to_execute_binance_futures_trade.invalid_trade_side", map[string]string{
				"side": trade.TradeSide,
			})
		}

		// Parse time in force.
		var timeInForce binanceproto.BinanceTimeInForce
		if orderType == binanceproto.BinanceOrderType_LIMIT {
			timeInForce = binanceproto.BinanceTimeInForce_GTC
		}

		orders = append(orders, &binanceproto.PerpetualFuturesOrder{
			Price:        float32(entry),
			OrderType:    orderType,
			Side:         side,
			Quantity:     float32(notionalSizes[i]),
			Symbol:       fmt.Sprintf("%s%s", trade.Asset, trade.Pair),
			TimeInForce:  timeInForce,
			PositionSide: binanceproto.BinancePositionSide_BOTH,
		})
	}

	// Add take profit orders to trade.
	for _, tp := range trade.TakeProfits {
		orders = append(orders, &binanceproto.PerpetualFuturesOrder{
			StopPrice:     float32(tp),
			OrderType:     binanceproto.BinanceOrderType_TAKE_PROFIT_MARKET,
			ClosePosition: true,
		})
	}

	// Execute trade.
	rsp, err := (&binanceproto.ExecuteFuturesPerpetualsTradeRequest{
		Orders:      orders,
		Timestamp:   timestamppb.Now(),
		Credentials: credentials,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_binance_futures_trade.binance", errParams)
	}

	return &FuturesTradeResponse{
		NotionalSize:           totalNotionalSize,
		ExchangeTradeID:        rsp.ExchangeTradeId,
		NumberOfExecutedOrders: int(rsp.NumberOfOrdersExecuted),
	}, nil
}

func sum(vs []float64) float64 {
	if len(vs) == 0 {
		return 0
	}

	return vs[0] + sum(vs[1:])
}
