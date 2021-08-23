package dao

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/domain"
	"time"

	"github.com/imdario/mergo"
	"github.com/monzo/terrors"
)

// ReadExchangeByExchangeID ...
func ReadExchangeByExchangeID(ctx context.Context, exchangeID string) (*domain.Exchange, error) {
	var (
		sql = `
		SELECT * FROM s_account_exchanges
		WHERE exchange_id=$1
		`
		exchanges []*domain.Exchange
	)

	if err := db.Select(ctx, exchanges, sql, exchangeID); err != nil {
		return nil, terrors.Propagate(err)
	}

	if len(exchanges) == 0 {
		return nil, terrors.NotFound("exchange-does-not-exist", "Failed to find exchange with exchange id", nil)
	}

	return exchanges[0], nil
}

// ListExchangesByUserID ...
func ListExchangesByUserID(ctx context.Context, userID string, isActive bool) ([]*domain.Exchange, error) {
	var (
		sql = `
		SELECT * FROM s_account_exchanges
		WHERE user_id=$1
		AND is_active=$2
		ORDER BY exchange
		`
		exchanges []*domain.Exchange
	)

	if err := db.Select(ctx, exchanges, sql, userID, isActive); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	if len(exchanges) == 0 {
		return nil, gerrors.NotFound("exchanges_not_found_for_user_id", nil)
	}

	return exchanges, nil
}

// AddExchange ...
func AddExchange(ctx context.Context, exchange *domain.Exchange) error {
	var (
		sql = `
		INSERT INTO s_account_exchanges
		(exchange, api_key, secret_key, user_id, created, updated)
		VALUES($1, $2, $3, $4, $5, $6)
		`
	)

	now := time.Now().UTC()
	if _, err := (db.Exec(
		ctx, sql,
		exchange.Exchange, exchange.APIKey, exchange.SecretKey, exchange.UserID,
		now, now,
	)); err != nil {
		return terrors.Propagate(err)
	}

	return nil
}

// RemoveExchange ...
func RemoveExchange(ctx context.Context, exchangeID string) error {
	var (
		sql = `
		DELETE FROM s_account_exchanges
		WHERE exchange_id=$1 
		`
	)

	if _, err := (db.Exec(
		ctx, sql, exchangeID,
	)); err != nil {
		return terrors.Propagate(err)
	}
	return nil
}

// UpdateExchange ...
func UpdateExchange(ctx context.Context, mutation *domain.Exchange) (*domain.Exchange, error) {
	var (
		sql = `
		UPDATE s_account_exchanges
		SET
		exchange=$1, api_key=$2, secret_key=$3, updated=$4
		`
	)
	if mutation.ExchangeID == "" {
		return nil, terrors.PreconditionFailed("missing-exchange-id", "Cannot update exchange with missing exchange id", nil)
	}

	exchange, err := ReadExchangeByExchangeID(ctx, mutation.ExchangeID)
	if err != nil {
		return nil, terrors.Propagate(err)
	}

	if err := mergo.Merge(&exchange, mutation); err != nil {
		return nil, terrors.BadRequest("mutation-merge-failure", "Failed to merge exchange mutation", map[string]string{
			"upstream_err": err.Error(),
		})
	}

	exchange.Updated = time.Now().UTC()

	if _, err := (db.Exec(
		ctx, sql,
		exchange.Exchange, exchange.APIKey, exchange.SecretKey, exchange.Updated,
	)); err != nil {
		return nil, terrors.Propagate(err)
	}

	return nil, nil
}
