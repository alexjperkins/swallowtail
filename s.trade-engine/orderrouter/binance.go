package orderrouter

import (
	"context"
	"fmt"
	"strings"

	"github.com/monzo/slog"
	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/risk"
	riskutil "swallowtail/libraries/risk"
	accountproto "swallowtail/s.account/proto"
	binanceproto "swallowtail/s.binance/proto"
	"swallowtail/s.trade-engine/domain"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

const (
	defaultBinanceDCAOrders = 5
)

func executeBinanceFuturesTrade(
	ctx context.Context,
	trade *domain.Trade,
	participant *tradeengineproto.AddParticipantToTradeRequest,
	credentials *binanceproto.Credentials,
) (*FuturesTradeResponse, error) {
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

	var positions []*risk.Position
	switch {
	case participant.Risk != 0:
		// Marshal side back into proto side.
		var side tradeengineproto.TRADE_SIDE
		switch strings.ToLower(trade.TradeSide) {
		case "long":
			side = tradeengineproto.TRADE_SIDE_LONG
		case "buy":
			side = tradeengineproto.TRADE_SIDE_BUY
		case "short":
			side = tradeengineproto.TRADE_SIDE_SHORT
		case "sell":
			side = tradeengineproto.TRADE_SIDE_SELL
		default:
			return nil, gerrors.Unimplemented("failed_to_execute_binance_futures_trade.unimplemented_trade_side", map[string]string{
				"trade_side": trade.TradeSide,
			})
		}

		// Marshal default dca strategy.
		var dcaStrategy tradeengineproto.DCA_STRATEGY
		switch account.DefaultDcaStrategy {
		case tradeengineproto.DCA_STRATEGY_CONSTANT.String():
			dcaStrategy = tradeengineproto.DCA_STRATEGY_CONSTANT
		case tradeengineproto.DCA_STRATEGY_LINEAR.String():
			dcaStrategy = tradeengineproto.DCA_STRATEGY_LINEAR
		case tradeengineproto.DCA_STRATEGY_EXPONENTIAL.String():
			dcaStrategy = tradeengineproto.DCA_STRATEGY_EXPONENTIAL
		}

		var entries = trade.Entries
		switch trade.OrderType {
		case tradeengineproto.ORDER_TYPE_DCA_FIRST_MARKET_REST_LIMIT.String(), tradeengineproto.ORDER_TYPE_MARKET.String():
			instrument := fmt.Sprintf("%s%s", trade.Asset, trade.Pair)
			currentPrice, err := getLatestPriceFromBinance(ctx, instrument)
			switch {
			case err != nil:
				slog.Error(ctx, "Failed to get the latest price for %s on binance; using all entries passed", instrument)
			default:
				// Set the entry for market order to the latest price so we can get the best risk possible.
				entries[len(entries)-1] = currentPrice
			}
		}

		// Calculate positins by risk & add as order.
		ps, err := riskutil.CalculatePositionsByRisk(trade.Entries, trade.StopLoss, float64(participant.Risk), defaultBinanceDCAOrders, side, dcaStrategy)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_calculate_risk_sizes.failed_to_calculate_notional_size_from_risk", nil)
		}
		positions = append(positions, ps...)
	default:
		return nil, gerrors.Unimplemented("failed_to_execute_binance_futures_trade.notional_size_calc_unimplimented", nil)
	}

	totalRisk := sumPositionsRisk(positions) * binanceAccountSize
	orders := make([]*binanceproto.PerpetualFuturesOrder, 0, len(trade.Entries)+1)

	// Add stop loss order. We add this first for safety.
	switch trade.StopLoss {
	case 0:
		slog.Warn(ctx, "Creating trade without a stop loss: [%v] %s", trade.ActorType, trade.ActorID)
	default:
		var side binanceproto.BinanceTradeSide

		switch trade.TradeSide {
		case tradeengineproto.TRADE_SIDE_BUY.String(), tradeengineproto.TRADE_SIDE_LONG.String():
			side = binanceproto.BinanceTradeSide_BINANCE_SELL
		case tradeengineproto.TRADE_SIDE_SELL.String(), tradeengineproto.TRADE_SIDE_SHORT.String():
			side = binanceproto.BinanceTradeSide_BINANCE_BUY
		}

		orders = append(orders, &binanceproto.PerpetualFuturesOrder{
			Side:          side,
			Symbol:        fmt.Sprintf("%s%s", trade.Asset, trade.Pair),
			StopPrice:     float32(trade.StopLoss),
			OrderType:     binanceproto.BinanceOrderType_BINANCE_STOP_MARKET,
			ClosePosition: true,
		})
	}

	errParams := map[string]string{
		"total_notional_size":   fmt.Sprintf("%v", totalRisk*binanceAccountSize),
		"risk":                  fmt.Sprintf("%v", participant.Risk),
		"total_risk_of_account": fmt.Sprintf("%v", totalRisk),
	}

	// Add positions as order.
	for i, riskedPosition := range positions {
		// Parse order type.
		var orderType binanceproto.BinanceOrderType
		switch {
		case trade.OrderType == tradeengineproto.ORDER_TYPE_MARKET.String():
			orderType = binanceproto.BinanceOrderType_BINANCE_MARKET
		case trade.OrderType == tradeengineproto.ORDER_TYPE_LIMIT.String():
			orderType = binanceproto.BinanceOrderType_BINANCE_LIMIT
		case trade.OrderType == tradeengineproto.ORDER_TYPE_DCA_ALL_LIMIT.String():
			orderType = binanceproto.BinanceOrderType_BINANCE_LIMIT
		case trade.OrderType == tradeengineproto.ORDER_TYPE_DCA_FIRST_MARKET_REST_LIMIT.String() && i == len(positions)-1:
			orderType = binanceproto.BinanceOrderType_BINANCE_MARKET
		case trade.OrderType == tradeengineproto.ORDER_TYPE_DCA_FIRST_MARKET_REST_LIMIT.String():
			orderType = binanceproto.BinanceOrderType_BINANCE_LIMIT
		default:
			slog.Warn(ctx, "Binance order router recieved trade with unrecognised order type: %s", trade.OrderType)
			orderType = binanceproto.BinanceOrderType_BINANCE_LIMIT
		}

		// Parse trade type.
		var side binanceproto.BinanceTradeSide
		switch trade.TradeSide {
		case tradeengineproto.TRADE_SIDE_BUY.String(), tradeengineproto.TRADE_SIDE_LONG.String():
			side = binanceproto.BinanceTradeSide_BINANCE_BUY
		case tradeengineproto.TRADE_SIDE_SELL.String(), tradeengineproto.TRADE_SIDE_SHORT.String():
			side = binanceproto.BinanceTradeSide_BINANCE_SELL
		default:
			return nil, gerrors.FailedPrecondition("failed_to_execute_binance_futures_trade.invalid_trade_side", map[string]string{
				"side": trade.TradeSide,
			})
		}

		// Parse time in force. Default is not required.
		var timeInForce binanceproto.BinanceTimeInForce
		if orderType == binanceproto.BinanceOrderType_BINANCE_LIMIT {
			timeInForce = binanceproto.BinanceTimeInForce_BINANCE_GTC
		}

		// We only set the price to a nonzero value if we have the determined the order type to be `LIMIT`.
		var price float32
		if orderType == binanceproto.BinanceOrderType_BINANCE_LIMIT {
			price = float32(riskedPosition.Price)
		}

		orders = append(orders, &binanceproto.PerpetualFuturesOrder{
			Price:        price,
			OrderType:    orderType,
			Side:         side,
			Quantity:     float32(riskedPosition.Risk * binanceAccountSize),
			Symbol:       fmt.Sprintf("%s%s", trade.Asset, trade.Pair),
			TimeInForce:  timeInForce,
			PositionSide: binanceproto.BinancePositionSide_BINANCE_SIDE_BOTH,
		})
	}

	// Add take profit orders to trade.
	// TODO; calculate risk.

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
		NotionalSize:           totalRisk * binanceAccountSize,
		ExchangeTradeID:        rsp.ExchangeTradeId,
		NumberOfExecutedOrders: int(rsp.NumberOfOrdersExecuted),
		ExecutionAlgoStrategy:  account.DefaultDcaStrategy,
	}, nil
}

func sumPositionsRisk(vs []*risk.Position) float64 {
	if len(vs) == 0 {
		return 0
	}

	return vs[0].Risk + sumPositionsRisk(vs[1:])
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
