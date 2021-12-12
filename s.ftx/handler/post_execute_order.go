package handler

import (
	"context"
	"strconv"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.ftx/client"
	"swallowtail/s.ftx/marshaling"
	ftxproto "swallowtail/s.ftx/proto"
)

// ExecuteOrder executes a given order on FTX.
func (s *FTXService) ExecuteOrder(
	ctx context.Context, in *ftxproto.ExecuteOrderRequest,
) (*ftxproto.ExecuteOrderResponse, error) {
	// Validate request.
	switch {
	case in.Order == nil:
		return nil, gerrors.BadParam("missing_param.orders", nil)
	case in.Credentials == nil:
		return nil, gerrors.BadParam("missing_param.credentials", nil)
	}

	// Validate credentials.
	if err := validateCredentials(in.Credentials); err != nil {
		return nil, gerrors.Unauthenticated("unauthorized.invalid_credentials", map[string]string{
			"msg": err.Error(),
		})
	}

	order := in.Order

	errParams := map[string]string{
		"actor_id":   order.ActorId,
		"instrument": order.Instrument,
		"asset":      order.Asset,
		"pair":       order.Pair.String(),
	}

	// Validate order.
	if err := validateOrder(in.Order); err != nil {
		slog.Error(ctx, "Failed to execute order invalid: Error: %v, Order: %+v", err, order)
		return nil, gerrors.Augment(err, "failed_to_execute_order.invalid_order", errParams)
	}

	// Marshal order to dto.
	dtoOrder, err := marshaling.OrderProtoToDTO(order)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_order", errParams)
	}

	// Marshal credentials.
	dtoCredentials := marshaling.VenueCredentialsProtoToFTXCredentials(in.Credentials)

	// Execute order via client.
	rsp, err := client.ExecuteOrder(ctx, dtoOrder, dtoCredentials)
	if err != nil {
		slog.Error(ctx, "Failed to execute order: Error: %v, Order: %+v", err, order)
		return nil, gerrors.Augment(err, "failed_to_execute_order", errParams)
	}

	// Embelish order with execution metadata.
	order.ExternalOrderId = strconv.Itoa(rsp.Result.ID)
	order.ExecutionTimestamp = rsp.Result.CreatedAt.UnixMilli()

	slog.Info(ctx, "Executed order: %+v", order)

	return &ftxproto.ExecuteOrderResponse{
		Order: order,
	}, nil
}
