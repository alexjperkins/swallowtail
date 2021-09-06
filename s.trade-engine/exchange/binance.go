package exchange

import (
	"context"
	"swallowtail/libraries/gerrors"
	binanceproto "swallowtail/s.binance/proto"
	"swallowtail/s.trade-engine/domain"
	tradeengineproto "swallowtail/s.trade-engine/proto"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func executeBinanceTrade(ctx context.Context, trade *domain.Trade, requestTimestamp time.Time) error {
	errParams := map[string]string{
		"trade_type": trade.TradeType,
	}

	tps := make([]float32, 0, len(trade.TakeProfits))
	for i, tp := range trade.TakeProfits {
		tps[i] = float32(tp)
	}

	switch trade.TradeType {
	case tradeengineproto.TRADE_TYPE_FUTURES_PERPETUALS.String():
		_, err := (&binanceproto.ExecuteFuturesPerpetualsTradeRequest{
			TradeId:     trade.ID,
			ActorId:     trade.ActorID,
			Side:        binanceproto.TradeSide_BUY,
			Asset:       trade.Asset,
			Pair:        trade.Pair,
			Entry:       float32(trade.Entry),
			StopLoss:    float32(trade.StopLoss),
			TakeProfits: tps,
			Timestamp:   timestamppb.New(requestTimestamp),
		}).Send(ctx).Response()
		if err != nil {
			return gerrors.Augment(err, "failed_to_execute_futures_perpetuals_trade_binance", errParams)
		}
	default:
		return gerrors.Unimplemented("failed_to_execute_binance_trade.unimplemented", errParams)
	}

	return nil
}
