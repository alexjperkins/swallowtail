package dao

import (
	"context"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/domain"
)

//  AddParticpantToTradeStrategy ...
func AddParticpantToTradeStrategy(ctx context.Context, tradeParticipant *domain.TradeStrategyParticipant) error {
	var (
		sql = `
		INSERT INTO 
			s_tradeengine_trade_strategy_participants(
				trade_strategy_id,
				user_id,
				is_bot,
				size,
				risk,
				venue,
				exchange_order_ids,
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
		tradeParticipant.TradeStrategyID, tradeParticipant.UserID, tradeParticipant.IsBot, tradeParticipant.Size, tradeParticipant.Risk, tradeParticipant.Venue, tradeParticipant.ExchangeOrderIDs,
		tradeParticipant.ExecutedTimestamp,
	)); err != nil {
		return gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	return nil
}

// ReadTradeStrategyParticipantByTradeStrategyID ...
func ReadTradeStrategyParticipantByTradeStrategyID(ctx context.Context, tradeStrategyID, userID string) (*domain.TradeStrategyParticipant, error) {
	var (
		sql = `
		SELECT * FROM s_tradeengine_trade_strategy_participants
		WHERE trade_strategy_id=$1
		AND user_id=$2
		`
		participants []*domain.TradeStrategyParticipant
	)

	if err := db.Select(ctx, &participants, sql, tradeStrategyID, userID); err != nil {
		return nil, gerrors.Propagate(err, gerrors.ErrUnknown, nil)
	}

	switch len(participants) {
	case 0:
		return nil, gerrors.NotFound("not_found.trade_strategy_participant", nil)
	case 1:
		return participants[0], nil
	default:
		slog.Critical(ctx, "Inconsistent data; more than one trade strategy participant with the same user id")
		return participants[0], nil
	}
}
