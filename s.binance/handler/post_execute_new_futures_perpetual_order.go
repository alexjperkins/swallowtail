package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	"swallowtail/s.binance/marshaling"
	binanceproto "swallowtail/s.binance/proto"
)

// ExecuteNewFuturesPerpetualOrder ...
func (s *BinanceService) ExecuteNewFuturesPerpetualOrder(
	ctx context.Context, in *binanceproto.ExecuteNewFuturesPerpetualOrderRequest,
) (*binanceproto.ExecuteNewFuturesPerpetualOrderResponse, error) {
	switch {
	case len(in.GetOrders()) == 0:
		return nil, gerrors.BadParam("missing_param.orders", nil)
	}

	// Validate credentials.
	if err := isValidCredentials(in.Credentials, false); err != nil {
		return nil, gerrors.Unauthenticated("invalid_credentials", nil)
	}

	errParams := map[string]string{
		"num_orders": strconv.Itoa(len(in.Orders)),
	}

	// Validate the trade.
	if err := validatePerpetualFuturesTrade(in); err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_perpetuals_trade.invalid_trade", errParams)
	}

	// Marshal orders.
	orders, err := marshaling.ProtoOrdersToExecutePerpetualsFutureTradeRequest(in.Orders)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_perpetuals_trade.invalid_order.marshaling", errParams)
	}

	// Marshal credentials.
	dtoCredentials := marshaling.CredentialsProtoToDTO(in.Credentials)

	// Execute orders synchronously.
	var (
		exchangeOrderIDs       strings.Builder
		maxTs                  int
		numberOfOrdersExecuted int64
	)
	for _, order := range orders {
		rsp, err := client.ExecutePerpetualFuturesTrade(ctx, order, dtoCredentials)
		if err != nil {
			isStopLoss := order.OrderType == binanceproto.BinanceOrderType_BINANCE_STOP.String() || order.OrderType == binanceproto.BinanceOrderType_BINANCE_STOP_MARKET.String()
			return nil, gerrors.Augment(err, "failed_to_execute_perpetuals_trade.order", map[string]string{
				"is_stop_loss": strconv.FormatBool(isStopLoss),
			})
		}

		exchangeOrderIDs.WriteString(fmt.Sprintf("%v,", rsp.OrderID))
		maxTs = max(maxTs, rsp.ExecutionTimestamp)
		numberOfOrdersExecuted++
	}

	return &binanceproto.ExecuteFuturesPerpetualsTradeResponse{
		ExchangeTradeId:        exchangeOrderIDs.String(),
		Timestamp:              int64(maxTs),
		NumberOfOrdersExecuted: numberOfOrdersExecuted,
	}, nil
}
