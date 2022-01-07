package orderrouter

import (
	"context"
	"swallowtail/libraries/gerrors"
	ftxproto "swallowtail/s.ftx/proto"

	"google.golang.org/protobuf/types/known/timestamppb"

	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// executeFTXNewOrder ...
func executeFTXNewOrders(ctx context.Context, order *tradeengineproto.Order, credentials *tradeengineproto.VenueCredentials) (*tradeengineproto.Order, error) {
	rsp, err := (&ftxproto.ExecuteNewOrderRequest{
		Order:       order,
		Credentials: credentials,
		Timestamp:   timestamppb.Now(),
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_route_and_execute_order.ftx", nil)
	}

	return rsp.Order, nil
}
