package handler

import (
	"context"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	"swallowtail/s.binance/marshaling"
	binanceproto "swallowtail/s.binance/proto"
)

// ExecuteNewSpotOrder ...
func (s *BinanceService) ExecuteNewSpotOrder(
	ctx context.Context, in *binanceproto.ExecuteNewSpotOrderRequest,
) (*binanceproto.ExecuteNewSpotOrderResponse, error) {
	// Validate request.
	switch {
	case in.Order == nil:
		return nil, gerrors.BadParam("missing_param.order", nil)
	case in.Credentials == nil:
		return nil, gerrors.BadParam("missing_param.credentials", nil)
	}

	// Validate credentials.
	if err := isValidCredentials(in.Credentials, false); err != nil {
		return nil, gerrors.Unauthenticated("invalid_credentials", nil)
	}

	order := in.GetOrder()

	errParams := map[string]string{
		"actor_id":   order.ActorId,
		"instrument": order.Instrument,
		"asset":      order.Asset,
		"pair":       order.Pair.String(),
	}

	// Validate order.
	if err := validateSpotOrder(order); err != nil {
		slog.Error(ctx, "Failed to validate order: Error: %v, Order: %+v", err, order)
		return nil, gerrors.Augment(err, "failed_to_execute_new_spot_order.invalid_order", errParams)
	}

	// Marshal order.
	dtoOrder, err := marshaling.ProtoOrderToExecuteSpotOrderRequest(order)
	if err != nil {
		slog.Error(ctx, "Failed to marshal order proto, Error: %v, Order: %+v", err, order)
		return nil, gerrors.Augment(err, "failed_to_execute_new_spot_order.marshal", errParams)
	}

	// Marshal credentials.
	dtoCredentials := marshaling.CredentialsProtoToDTO(in.Credentials)

	// Execute order.
	rsp, err := client.ExecuteSpotOrder(ctx, dtoOrder, dtoCredentials)
	if err != nil {
		slog.Error(ctx, "Failed to execute spot order, Error: %v, Order: %+v", err, order)
		return nil, gerrors.Augment(err, "failed_to_execute_new_spot_order.order", nil)
	}

	// Embelish order with exchange metadata.
	order.ExternalOrderId = rsp.ExternalOrderID
	order.ExecutionTimestamp = int64(rsp.ExecutionTimestamp)

	slog.Info(ctx, "Order executed: %+v", order)

	return nil, gerrors.Unimplemented("failed_to_execute_new_spot_order.unimplemented", nil)
}
