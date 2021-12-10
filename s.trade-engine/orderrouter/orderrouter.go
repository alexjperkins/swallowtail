package orderrouter

import (
	"context"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/libraries/gerrors"
	binanceproto "swallowtail/s.binance/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// RouteAndExecuteNewOrder routes order to the correct exchange & executes.
func RouteAndExecuteNewOrder(
	ctx context.Context,
	order *tradeengineproto.Order,
	venue tradeengineproto.VENUE,
	instrumentType tradeengineproto.INSTRUMENT_TYPE,
	venueCredentials *tradeengineproto.VenueCredentials,
) (*tradeengineproto.Order, error) {
	errParams := map[string]string{
		"venue_id":        strings.ToLower(venue.String()),
		"instrument_type": strings.ToLower(instrumentType.String()),
	}

	switch venue {
	case tradeengineproto.VENUE_BINANCE:
		switch instrumentType {
		case tradeengineproto.INSTRUMENT_TYPE_FUTURE_PERPETUAL:
			return executeBinanceNewPerpetualFuturesOrders(ctx, order, venueCredentials)
		case tradeengineproto.INSTRUMENT_TYPE_SPOT:
			return executeBinanceNewSpotOrders(ctx, order, venueCredentials)
		default:
			return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.instrument_exchange_pair_umimplemented", errParams)
		}
	case tradeengineproto.VENUE_FTX:
		return executeFTXNewOrders(ctx, order, venueCredentials)
	default:
		return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.exchange_unimplemented", errParams)
	}
}

// executeBinanceNewPerpetualFuturesOrder ...
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

// executeBinanceNewSpotOrder ...
func executeBinanceNewSpotOrders(ctx context.Context, order *tradeengineproto.Order, credentials *tradeengineproto.VenueCredentials) (*tradeengineproto.Order, error) {
	return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.ftx_new_orders_spot_unimplemented", nil)
}

// executeFTXNewOrder ...
func executeFTXNewOrders(ctx context.Context, order *tradeengineproto.Order, credentials *tradeengineproto.VenueCredentials) (*tradeengineproto.Order, error) {
	return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.ftx_new_orders_unimplemented", nil)
}
