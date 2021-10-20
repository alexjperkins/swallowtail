package orderrouter

import (
	"context"

	"swallowtail/libraries/gerrors"
	binanceproto "swallowtail/s.binance/proto"
)

func getLatestPriceFromBinance(ctx context.Context, instrument string) (float64, error) {
	rsp, err := (&binanceproto.GetLatestPriceRequest{
		Symbol: instrument,
	}).Send(ctx).Response()
	if err != nil {
		return 0, gerrors.Augment(err, "failed_to_get_latest_price", nil)
	}

	return float64(rsp.Price), nil
}
