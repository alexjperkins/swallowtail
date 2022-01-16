package orderrouter

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/libraries/gerrors"
	binanceproto "swallowtail/s.binance/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// executeBinanceNewPerpetualFuturesOrder executes an order as a binance perpetual future order.
func executeBinanceNewPerpetualFuturesOrders(ctx context.Context, order *tradeengineproto.Order, credentials *tradeengineproto.VenueCredentials) (*tradeengineproto.Order, error) {
	rsp, err := (&binanceproto.ExecuteNewFuturesPerpetualOrderRequest{
		Order:       order,
		Credentials: credentials,
		Timestamp:   timestamppb.Now(),
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_route_and_execute_order.binance_perpetual_futures", nil)
	}

	return rsp.Order, nil
}

// executeBinanceNewSpotOrder executes an order as a binance spot order.
func executeBinanceNewSpotOrders(ctx context.Context, order *tradeengineproto.Order, credentials *tradeengineproto.VenueCredentials) (*tradeengineproto.Order, error) {
	return nil, gerrors.Unimplemented("unimplemented.binance_spot_orders", nil)
}
