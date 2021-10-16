package handler

import (
	"context"
	"strconv"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.ftx/marshaling"
	ftxproto "swallowtail/s.ftx/proto"
)

func (s *FTXService) ExecuteOrder(
	ctx context.Context, in *ftxproto.ExecuteOrderRequest,
) (*ftxproto.ExecuteOrderResponse, error) {
	switch {
	case len(in.Orders) == 0:
		return nil, gerrors.BadParam("missing_param.orders", nil)
	}

	// Validate credentials.
	if err := validateCredentials(in.Credentials); err != nil {
		return nil, gerrors.Unauthenticated("unauthorized.invalid_credentials", map[string]string{
			"msg": err.Error(),
		})
	}

	errParams := map[string]string{
		"total_num_orders": strconv.Itoa(len(in.Orders)),
	}

	// Validate all orders.
	if err := validateOrders(in.Orders); err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_order.invalid_order", errParams)
	}

	orders, err := marshaling.OrdersProtoToDTO(in.Orders)

	return nil, gerrors.Unimplemented("failed_to_execute_order.unimplemented", nil)
}
