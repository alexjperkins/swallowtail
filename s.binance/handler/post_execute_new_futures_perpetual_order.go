package handler

import (
	"context"
	"strconv"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	"swallowtail/s.binance/marshaling"
	binanceproto "swallowtail/s.binance/proto"
)

// ExecuteNewFuturesPerpetualOrder excutes a futures perpetual order on binance.
func (s *BinanceService) ExecuteNewFuturesPerpetualOrder(
	ctx context.Context, in *binanceproto.ExecuteNewFuturesPerpetualOrderRequest,
) (*binanceproto.ExecuteNewFuturesPerpetualOrderResponse, error) {
	switch {
	case in.Order == nil:
		return nil, gerrors.BadParam("missing_param.order", nil)
	case in.Credentials == nil:
		return nil, gerrors.BadParam("missing_param.credentials", nil)
	}

	order := in.GetOrder()

	// Validate credentials.
	if err := isValidCredentials(in.Credentials, false); err != nil {
		return nil, gerrors.Unauthenticated("invalid_credentials", nil)
	}

	// Validate the order.
	if err := validatePerpetualFuturesOrder(order); err != nil {
		slog.Error(ctx, "Invalid order: Error: %v, Order: %+v", err, order)
		return nil, gerrors.Augment(err, "failed_to_execute_perpetuals_trade.invalid_order", nil)
	}

	// Marshal orders.
	dtoOrder, err := marshaling.ProtoOrderToExecutePerpetualsFutureOrderRequest(order)
	if err != nil {
		slog.Error(ctx, "Failed to marshal order: Error: %v, Order: %+v", err, order)
		return nil, gerrors.Augment(err, "failed_to_execute_perpetuals_trade.invalid_order.marshaling", nil)
	}

	// Marshal credentials.
	dtoCredentials := marshaling.CredentialsProtoToDTO(in.Credentials)

	// Execute order.
	rsp, err := client.ExecutePerpetualFuturesOrder(ctx, dtoOrder, dtoCredentials)
	if err != nil {
		slog.Error(ctx, "Failed to execute order: Error: %v, Order: %+v", err, rsp)
		return nil, gerrors.Augment(err, "failed_to_execute_perpetuals_trade.order", nil)
	}

	// Update order with metadata from exchange.
	order.ExternalOrderId = strconv.Itoa(rsp.OrderID)
	order.ExecutionTimestamp = int64(rsp.ExecutionTimestamp)

	slog.Info(ctx, "Order executed: %+v", order)

	return &binanceproto.ExecuteNewFuturesPerpetualOrderResponse{
		Order: order,
	}, nil
}
