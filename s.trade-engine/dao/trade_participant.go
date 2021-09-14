package dao

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/domain"

	"github.com/monzo/slog"
)

// AddTradeParticpantToTrade ...
func AddTradeParticpantToTrade(ctx context.Context, tradeParticipant *domain.TradeParticipant) error {
	var (
		sql = `
		INSERT INTO 
			s_tradeengine_trade_participants(
				trade_id,
				user_id,
				is_bot,
				size,
				risk,
				exchange,
				exchange_order_id,
				executed
			)
		VALUES
			(
				$1, $2, $3, $4, $5, $6, $7, $8
			)
		`
	)

	if _, err := (db.Exec(
		ctx, sql,
		tradeParticipant.TradeID, tradeParticipant.UserID, tradeParticipant.IsBot, tradeParticipant.Size, tradeParticipant.Risk, tradeParticipant.Exchange, tradeParticipant.ExchangeOrderID,
		tradeParticipant.ExecutedTimestamp,
	)); err != nil {
		return gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return nil
}

// ReadTradeParticipantByTradeID ...
func ReadTradeParticipantByTradeID(ctx context.Context, tradeID, userID string) (*domain.TradeParticipant, error) {
	var (
		sql = `
		SELECT * FROM s_tradeengine_trade_participants
		WHERE trade_id=$1
		AND user_id=$2
		`
		participants []*domain.TradeParticipant
	)

	if err := db.Select(ctx, &participants, sql, tradeID, userID); err != nil {
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
