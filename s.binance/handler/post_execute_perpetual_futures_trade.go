package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	"swallowtail/s.binance/exchangeinfo"
	"swallowtail/s.binance/marshaling"
	binanceproto "swallowtail/s.binance/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// ExecuteFuturesPerpetualsTrade ...
func (s *BinanceService) ExecuteFuturesPerpetualsTrade(
	ctx context.Context, in *binanceproto.ExecuteFuturesPerpetualsTradeRequest,
) (*binanceproto.ExecuteFuturesPerpetualsTradeResponse, error) {
	switch {
	case in.Asset == "":
		return nil, gerrors.BadParam("missing_param.asset", nil)
	case in.Pair == "":
		return nil, gerrors.BadParam("missing_param.pair", nil)
	}

	// Validate credentials.
	if err := isValidCredentials(in.Credentials, false); err != nil {
		return nil, gerrors.Unauthenticated("invalid_credentials", nil)
	}

	errParams := map[string]string{
		"asset": in.Asset,
		"pair":  in.Pair,
	}

	// Validate the trade.
	if err := validatePerpetualFuturesTrade(in); err != nil {
		return nil, gerrors.Augment(err, "failed_to_execute_perpetuals_trade.invalid_trade", errParams)
	}

	// Round the quantity to the minimum precision allowed on the exchange.
	assetQuantityPrecision, ok := exchangeinfo.GetBaseAssetQuantityPrecision(in.Asset)
	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_execute_perpetuals_trade.asset_quantity_precision_unknown", errParams)
	}

	// Round the price to the minimum precision allowed on the exchange.
	assetPricePrecision, ok := exchangeinfo.GetBaseAssetPricePrecision(in.Asset)
	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_execute_perpetuals_trade.asset_price_precision_unknown", errParams)
	}

	// Convert floats to minimum precision rounded strings.
	quantity := roundToPrecisionString(float64(in.NotionalSize), assetQuantityPrecision)
	entryPrice := roundToPrecisionString(float64(in.Entry), assetPricePrecision)
	stopLossPrice := roundToPrecisionString(float64(in.StopLoss), assetPricePrecision)

	orders := []*client.ExecutePerpetualFuturesTradeRequest{}

	// Add Stop Loss.
	orders = append(orders, &client.ExecutePerpetualFuturesTradeRequest{
		Symbol:           strings.ToUpper(fmt.Sprintf("%s%s", in.Asset, in.Pair)),
		StopPrice:        stopLossPrice,
		Side:             "SELL",
		Type:             tradeengineproto.ORDER_TYPE_STOP_MARKET.String(),
		ClosePosition:    "true",
		NewOrderRespType: "ACK",
		WorkingType:      "MARK_PRICE",
	})

	// Add Entry.
	entry := &client.ExecutePerpetualFuturesTradeRequest{
		Symbol:           strings.ToUpper(fmt.Sprintf("%s%s", in.Asset, in.Pair)),
		Side:             convertLongAndShort(in.TradeSide),
		Type:             strings.ToUpper(in.OrderType),
		Quantity:         quantity,
		NewOrderRespType: "ACK",
	}

	// Decorate entry if we have a LIMIT order.
	if strings.ToUpper(in.OrderType) == tradeengineproto.ORDER_TYPE_LIMIT.String() {
		entry.Price = entryPrice
		entry.TimeInForce = "GTC"
	}

	// Add entry to orders
	orders = append(orders, entry)

	// Marshal credentials
	dtoCredentials := marshaling.CredentialsProtoToDTO(in.Credentials)

	// Execute orders synchronously.
	var (
		exchangeID strings.Builder
		maxTs      int
	)
	for _, order := range orders {
		rsp, err := client.ExecutePerpetualFuturesTrade(ctx, order, dtoCredentials)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_execute_perpetuals_trade.order", map[string]string{
				// TODO: this should be improved somewhat.
				"is_stop_loss": strconv.FormatBool(order.ClosePosition == "true"),
				"is_entry":     strconv.FormatBool(order.Quantity != ""),
			})
		}

		exchangeID.WriteString(fmt.Sprintf("%v,", rsp.OrderID))
		maxTs = max(maxTs, rsp.ExecutionTimestamp)
	}

	return &binanceproto.ExecuteFuturesPerpetualsTradeResponse{
		ExchangeTradeId: exchangeID.String(),
		Timestamp:       int64(maxTs),
	}, nil
}
