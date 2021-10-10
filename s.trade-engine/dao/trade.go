package dao

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/domain"
	"time"

	"github.com/monzo/slog"
)

// TradeExists checks if the trade already exists in persistent storage.
func TradeExists(ctx context.Context, idempotencyKey string) (bool, error) {
	var (
		sql = `
		SELECT * FROM s_tradeengine_trades
		WHERE
			idempotency_key=$1
		`
		trades []*domain.Trade
	)

	err := db.Select(ctx, &trades, sql, idempotencyKey)
	if err != nil {
		return false, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return len(trades) > 0, nil
}

// CreateTrade ...
func CreateTrade(ctx context.Context, trade *domain.Trade) error {
	var (
		sql = `
		INSERT INTO 
			s_tradeengine_trades(
				actor_id,
				humanized_actor_name,
				actor_type,
				idempotency_key,
				order_type,
				trade_type,
				trade_side,
				asset,
				pair,
				entries,
				stop_loss,
				take_profits,
				current_price,
				status,
				created,
				last_updated
			)
		VALUES
			(
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
			)
		`
	)

	t := trade
	now := time.Now().UTC()
	t.Created = now
	t.LastUpdated = now

	if _, err := (db.Exec(
		ctx, sql,
		t.ActorID, t.HumanizedActorName, t.ActorType, t.IdempotencyKey, t.OrderType, t.TradeType, t.TradeSide, trade.Asset, t.Pair, t.Entries, t.StopLoss, t.TakeProfits, t.CurrentPrice, t.Status,
		t.Created, t.LastUpdated,
	)); err != nil {
		return gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return nil
}

// ReadTradeByTradeID ...
func ReadTradeByTradeID(ctx context.Context, tradeID string) (*domain.Trade, error) {
	var (
		sql = `
		SELECT * FROM s_tradeengine_trades
		WHERE trade_id=$1
		`
		trades []*domain.Trade
	)

	if err := db.Select(ctx, &trades, sql, tradeID); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(trades) {
	case 0:
		return nil, gerrors.NotFound("not_found.trade", nil)
	case 1:
		return trades[0], nil
	default:
		// This should never happen. But if it does we at least want a record of it.
		slog.Critical(ctx, "Critical State: more than one identical trade.", map[string]string{
			"trade_id": tradeID,
		})
		return trades[0], nil
	}
}

// ReadTradeByIdempotencyKey ...
func ReadTradeByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Trade, error) {
	var (
		sql = `
		SELECT * FROM s_tradeengine_trades
		WHERE idempotency_key=$1
		`
		trades []*domain.Trade
	)

	if err := db.Select(ctx, &trades, sql, idempotencyKey); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(trades) {
	case 0:
		return nil, gerrors.NotFound("not_found.trade", nil)
	case 1:
		return trades[0], nil
	default:
		// This should never happen. But if it does we at least want a record of it.
		slog.Critical(ctx, "Critical State: more than one identical trade.", map[string]string{
			"trade_id": idempotencyKey,
		})
		return trades[0], nil
	}
}
