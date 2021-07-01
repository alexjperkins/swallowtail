package handler

import (
	"context"
	"time"

	binanceclient "swallowtail/s.binance/client"
	"swallowtail/s.binance/dao"
	"swallowtail/s.binance/domain"
	binanceproto "swallowtail/s.binance/proto"

	"github.com/monzo/terrors"
)

func (b *BinanceService) handlePUTSpotTrade(
	ctx context.Context, in *binanceproto.SpotTradeRequest,
) (*binanceproto.SpotTradeResponse, error) {
	errParams := map[string]string{
		"discord_user_id":     in.UserDiscordId,
		"idempotency_key":     in.IdempotencyKey,
		"asset_pair":          in.AssetPair,
		"amount":              in.Amount,
		"value":               in.Value,
		"trade_type":          binanceproto.TradeType_SPOT.String(),
		"trade_side":          in.Side.String(),
		"attempt_retry_until": in.AttemptRetryUntil.String(),
	}

	trade, err := dao.Exists(ctx, in.IdempotencyKey)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to check if trade has already been made", errParams)
	}
	if trade != nil {
		return &binanceproto.SpotTradeResponse{
			Executed: false,
			TradeId:  trade.TradeID,
		}, nil
	}

	trade = &domain.Trade{
		UserDiscordID:  in.UserDiscordId,
		IdempotencyKey: in.IdempotencyKey,
		Side:           in.Side.String(),
		Type:           binanceproto.TradeType_SPOT.String(),
		AssetPair:      in.AssetPair,
		Amount:         in.Amount,
		Value:          in.Value,
		Created:        time.Now().UTC(),
	}

	trade, err = executeTradeWithRetry(ctx, binanceclient.ExecuteSpotTrade, trade, 5)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to execute trade", errParams)
	}

	if err := dao.SetTrade(ctx, trade); err != nil {
		return nil, terrors.Augment(err, "Failed to set trade; trade executed", errParams)
	}

	return &binanceproto.SpotTradeResponse{
		Executed: true,
	}, nil
}
