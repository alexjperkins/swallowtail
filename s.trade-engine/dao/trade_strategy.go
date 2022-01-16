package dao

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/domain"
	"time"

	"github.com/monzo/slog"
)

// TradeStrategyExists checks if the trade already exists in persistent storage.
func TradeStrategyExists(ctx context.Context, idempotencyKey string) (bool, error) {
	var (
		sql = `
		SELECT * FROM s_tradeengine_trade_strategies
		WHERE
			idempotency_key=$1
		`
		trades []*domain.TradeStrategy
	)

	err := db.Select(ctx, &trades, sql, idempotencyKey)
	if err != nil {
		return false, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return len(trades) > 0, nil
}

// CreateTradeStrategy ...
func CreateTradeStrategy(ctx context.Context, trade *domain.TradeStrategy) error {
	var (
		sql = `
		INSERT INTO 
			s_tradeengine_trade_strategies(
				actor_id,
				humanized_actor_name,
				actor_type,
				idempotency_key,
				execution_strategy,
				instrument_type,
				trade_side,
				asset,
				pair,
				entries,
				stop_loss,
				take_profits,
				current_price,
				status,
				tradeable_venues,
				created,
				last_updated
			)
		VALUES
			(
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
			)
		`
	)

	t := trade
	now := time.Now().UTC()
	t.Created = now
	t.LastUpdated = now

	if _, err := (db.Exec(
		ctx, sql,
		t.ActorID, t.HumanizedActorName, t.ActorType, t.IdempotencyKey, t.ExecutionStrategy, t.InstrumentType, t.TradeSide, trade.Asset, t.Pair, t.Entries, t.StopLoss, t.TakeProfits, t.CurrentPrice, t.Status, t.TradeableVenues,
		t.Created, t.LastUpdated,
	)); err != nil {
		return gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return nil
}

// ReadTradeStrategyByTradeStrategyID ...
func ReadTradeStrategyByTradeStrategyID(ctx context.Context, tradeID string) (*domain.TradeStrategy, error) {
	var (
		sql = `
		SELECT * FROM s_tradeengine_trade_strategies
		WHERE trade_strategy_id=$1
		`
		tradeStrategies []*domain.TradeStrategy
	)

	if err := db.Select(ctx, &tradeStrategies, sql, tradeID); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(tradeStrategies) {
	case 0:
		return nil, gerrors.NotFound("not_found.trade_strategy", nil)
	case 1:
		return tradeStrategies[0], nil
	default:
		// This should never happen. But if it does we at least want a record of it.
		slog.Critical(ctx, "Critical State: more than one identical trade strategy.", map[string]string{
			"trade_strategy_id": tradeID,
		})
		return tradeStrategies[0], nil
	}
}

// ReadTradeStrategyByIdempotencyKey ...
func ReadTradeStrategyByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.TradeStrategy, error) {
	var (
		sql = `
		SELECT * FROM s_tradeengine_trade_strategies
		WHERE idempotency_key=$1
		`
		tradeStrategies []*domain.TradeStrategy
	)

	if err := db.Select(ctx, &tradeStrategies, sql, idempotencyKey); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(tradeStrategies) {
	case 0:
		return nil, gerrors.NotFound("not_found.trade_strategy", nil)
	case 1:
		return tradeStrategies[0], nil
	default:
		// This should never happen. But if it does we at least want a record of it.
		slog.Critical(ctx, "Critical State: more than one identical trade.", map[string]string{
			"trade_strategy_id": idempotencyKey,
		})
		return tradeStrategies[0], nil
	}
}
