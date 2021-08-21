package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	coingeckoproto "swallowtail/s.coingecko/proto"
)

func getLatestPriceBySymbol(ctx context.Context, symbol, assetPair string) (*coingeckoproto.GetAssetLatestPriceBySymbolResponse, error) {
	rsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
		AssetSymbol: symbol,
		AssetPair:   assetPair,
	}).SendWithTimeout(ctx, 30*time.Second).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "fetch_latest_price_failed", map[string]string{
			"asset_symbol": symbol,
			"asset_pair":   assetPair,
		})
	}

	return rsp, nil
}

func pageAccount(ctx context.Context, userID, msg, spreadsheetID string) error {
	wrappedMsg := fmt.Sprintf("%s\n`spreadsheet url: %s`", msg, spreadsheetID)
	if _, err := (&accountproto.PageAccountRequest{
		UserId:   userID,
		Content:  wrappedMsg,
		Priority: accountproto.PagerPriority_HIGH,
	}).Send(ctx).Response(); err != nil {
		slog.Warn(ctx, "Failed to send: %v to %v", msg, userID)
		return gerrors.Augment(err, "Failed to page account", map[string]string{
			"user_id": userID,
		})
	}
	return nil
}
