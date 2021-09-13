package handler

import (
	"context"

	"github.com/monzo/slog"
)

func executeTradeForUser(ctx context.Context, userID, tradeID string, riskPercentage int) error {
	if err := notifyUserOnSuccess(ctx, userID, tradeID, riskPercentage); err != nil {
		slog.Error(ctx, "Failed to notify user of trade.", nil)
	}

	return nil
}
