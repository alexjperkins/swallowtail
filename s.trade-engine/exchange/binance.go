package exchange

import (
	"context"
	"fmt"
	"strings"

	"swallowtail/libraries/gerrors"
	riskutil "swallowtail/libraries/risk"
	binanceproto "swallowtail/s.binance/proto"
	"swallowtail/s.trade-engine/domain"
	tradeengineproto "swallowtail/s.trade-engine/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func calculuateRiskFromBinanceAccount(ctx context.Context, binanceCredentials *binanceproto.Credentials, entry, stopLoss, risk float64, side *tradeengineproto.TRADE_SIDE) (float64, error) {
	rsp, err := (&binanceproto.ReadPerpetualFuturesAccountRequest{
		Credentials: binanceCredentials,
		ActorId:     binanceproto.BinanceAccountActorTradeEngineSystem,
	}).Send(ctx).Response()
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_calculate_risk_size.failed_to_read_binance_account_futures_balance", nil)
	}

	accountSize := float64(rsp.Balance)
	notionalSize, err := riskutil.CalculateNotionalSizeFromPositionAndRisk(entry, stopLoss, risk, accountSize, side)
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_calculate_risk_size.failed_to_calculate_notional_size_from_risk", nil)
	}

	return notionalSize, nil
}

func executeBinanceFuturesTrade(
	ctx context.Context,
	trade *domain.Trade,
	participant *tradeengineproto.AddParticipantToTradeRequest,
	credentials *binanceproto.Credentials,
) (*FuturesTradeResponse, error) {

	var notionalSize float64
	switch {
	case participant.Size == 0:
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

		ns, err := calculuateRiskFromBinanceAccount(ctx, credentials, trade.Entry, trade.StopLoss, float64(participant.Risk), side)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_execute_binance_futures_trade", map[string]string{
				"side": side.String(),
			})
		}
		notionalSize = ns
	default:
		notionalSize = float64(participant.Size)
	}

	errParams := map[string]string{
		"notional_size": fmt.Sprintf("%v", notionalSize),
	}

	rsp, err := (&binanceproto.ExecuteFuturesPerpetualsTradeRequest{
		TradeSide:    trade.TradeSide,
		OrderType:    trade.OrderType,
		NotionalSize: float32(notionalSize),
		Asset:        trade.Asset,
		Pair:         trade.Pair,
		Entry:        float32(trade.Entry),
		StopLoss:     float32(trade.StopLoss),
		Timestamp:    timestamppb.Now(),
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_binance_futures_trade.binance", errParams)
	}

	return &FuturesTradeResponse{
		NotionalSize:    notionalSize,
		ExchangeTradeID: rsp.ExchangeTradeId,
	}, nil
}
