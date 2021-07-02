package dao

import (
	"context"
	"swallowtail/s.binance/domain"

	"github.com/monzo/terrors"
)

// Exists returns true if the given idempotency key already exists in the index; otherwise returns false.
func Exists(ctx context.Context, idempotencyKey string) (*domain.Trade, error) {
	var (
		sql = `
		SELECT * FROM s_binance_trades
		WHERE idempotency_key=$1
		`
		trades []*domain.Trade
	)
	if err := db.Select(ctx, &trades, sql, idempotencyKey); err != nil {
		return nil, terrors.Propagate(err)
	}

	if len(trades) == 0 {
		return nil, nil
	}
	// We enforce uniqueness of the idempotency key at the data level; so we should only expect the one
	// here if we do already have a trade exist with this idempotency key.
	return trades[0], nil
}

// SetTrade ...
func SetTrade(ctx context.Context, trade *domain.Trade) error {
	return nil
}
