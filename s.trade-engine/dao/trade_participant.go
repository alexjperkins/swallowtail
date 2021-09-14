package dao

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/domain"

	"github.com/monzo/slog"
)

// AddTradeParticpantToTrade ...
func AddTradeParticpantToTrade(ctx context.Context, tradeParticipant *domain.TradeParticipent) error {
	var (
		sql = `
		INSERT INTO 
			s_tradeengine_trade_participants(
			)
		VALUES
			(
			)
		`
	)

	if _, err := (db.Exec(ctx, sql)); err != nil {
		gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return nil
}

// ReadTradeParticipantByTradeID ...
func ReadTradeParticipantByTradeID(ctx context.Context, tradeID, userID string) (*domain.TradeParticipent, error) {
	var (
		sql = `
		SELECT * FROM s_tradeengine_trade_participants
		WHERE
			trade_id=$1
		AND
			user_id=$2
		`
		participants []*domain.TradeParticipent
	)

	if err := db.Select(ctx, sql, tradeID, userID); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(participants) {
	case 0:
		return nil, gerrors.NotFound("not_found.trade_participant", nil)
	case 1:
		return participants[0], nil
	default:
		slog.Critical(ctx, "Inconsistent data; more than one trade participant with the same user id")
		return participants[0], nil
	}
}
