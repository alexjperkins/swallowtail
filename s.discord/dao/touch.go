package dao

import (
	"context"
	"swallowtail/s.discord/domain"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/imdario/mergo"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

// Exists checkcks
func Exists(ctx context.Context, idempotencyKey string) (*domain.Touch, bool, error) {
	var (
		sql     = `SELECT * FROM s_discord_touches WHERE idempotency_key=$1`
		touches []*domain.Touch
	)
	if idempotencyKey == "" {
		return nil, false, nil
	}

	err := pgxscan.Select(ctx, db, &touches, sql, idempotencyKey)
	if err != nil {
		return nil, false, terrors.Propagate(err)
	}
	switch {
	case len(touches) > 1:
		slog.Critical(ctx, "We have more than one touch with identical idempotency key", map[string]string{
			"idempotency_key": idempotencyKey,
		})
	case len(touches) == 0:
		return nil, false, nil
	}
	return touches[0], true, nil
}

// Update updates existing touch via merging with mutation & persisting.
func Update(ctx context.Context, mutation *domain.Touch) (*domain.Touch, error) {
	var (
		sql = `
		UPDATE s_discord_touches
		SET idempotency_key=$1, udpated=$2, sender_id=$3`
	)

	if mutation.IdempotencyKey == "" {
		slog.Warn(ctx, "Attempting to update touch with nil idempotency key")
		return nil, nil
	}

	existing, exists, err := Exists(ctx, mutation.IdempotencyKey)
	if err != nil {
		return nil, terrors.Propagate(err)
	}

	switch {
	case exists:
		if err := mergo.Merge(existing, mutation); err != nil {
			return nil, terrors.BadRequest("mutation-merge-failure", "Failed to merge touch mutation", map[string]string{
				"upstream_err": err.Error(),
			})
		}

		if _, err := (db.Exec(
			ctx, sql, existing.IdempotencyKey, existing.Updated, existing.SenderID,
		)); err != nil {
			return nil, terrors.Propagate(err)
		}

		return existing, nil
	default:
		return nil, terrors.NotFound("touch-not-found", "Cannot update; touch with idempotency_key not found", map[string]string{
			"idempotency_key": mutation.IdempotencyKey,
		})
	}
}

func Create(ctx context.Context, touch *domain.Touch) (*domain.Touch, error) {
	var (
		sql = `
		INSERT INTO s_discord_touches
		(idempotency_key, updated, sender_id)
		VALUES
		($1, $2, $3)
		`
	)

	if touch.IdempotencyKey == "" {
		return nil, terrors.BadRequest("bad-touch", "Failed to create touch with nil idempotency key", nil)
	}

	if _, err := (db.Exec(
		ctx, sql, touch.IdempotencyKey, touch.Updated, touch.SenderID,
	)); err != nil {
		return nil, terrors.Propagate(err)
	}

	return touch, nil
}
