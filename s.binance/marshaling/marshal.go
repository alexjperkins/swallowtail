package marshaling

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	"swallowtail/s.binance/exchangeinfo"
	binanceproto "swallowtail/s.binance/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// CredentialsProtoToDTO ...
func CredentialsProtoToDTO(in *tradeengineproto.VenueCredentials) *client.Credentials {
	return &client.Credentials{
		APIKey:    in.ApiKey,
		SecretKey: in.SecretKey,
	}
}

// VerifyRequestDTOToProto ...
func VerifyRequestDTOToProto(in *client.VerifyCredentialsResponse) *binanceproto.VerifyCredentialsResponse {
	isSuccess, reason := isSuccess(in)

	return &binanceproto.VerifyCredentialsResponse{
		Success:         isSuccess,
		ReadEnabled:     in.EnableReading,
		FuturesEnabled:  in.EnableFutures,
		WithdrawEnabled: in.EnableWithdrawals,
		SpotEnabled:     in.EnableSpotAndMarginTrading,
		OptionsEnabled:  in.EnableVanillaOptions,
		IpRestrictions:  in.IPRestrict,
		Reason:          reason,
	}
}

// PerpetualFuturesAccountBalanceDTOToProto ...
func PerpetualFuturesAccountBalanceDTOToProto(in *client.PerpetualFuturesAccountBalance) (*binanceproto.ReadPerpetualFuturesAccountResponse, error) {
	balance, err := strconv.ParseFloat(in.Balance, 64)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_float.balance", nil)
	}

	availableBalance, err := strconv.ParseFloat(in.AvailableBalance, 64)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_float.available_balance", nil)
	}

	return &binanceproto.ReadPerpetualFuturesAccountResponse{
		Asset:            in.Asset,
		Balance:          float32(balance),
		AvailableBalance: float32(availableBalance),
		LastUpdated:      timestamppb.New(time.Unix(int64(in.LastUpdated/1_000), 0)),
	}, nil
}

// ProtoOrderToExecutePerpetualsFutureOrderRequest ...
func ProtoOrderToExecutePerpetualsFutureOrderRequest(in *tradeengineproto.Order) (*client.ExecutePerpetualFuturesOrderRequest, error) {
	// Parse symbol.
	var symbol string
	switch {
	case in.Asset == "" && in.Instrument == "":
		return nil, gerrors.FailedPrecondition("missing_param.instrument_or_asset", nil)
	case in.Instrument == "":
		symbol = fmt.Sprintf("%s%s", strings.ToUpper(in.Asset), strings.ToUpper(in.Pair.String()))
	default:
		symbol = strings.ToUpper(in.Instrument)
	}

	errParams := map[string]string{
		"symbol": symbol,
	}

	// Round the quantity to the minimum precision allowed on the exchange.
	assetQuantityPrecision, ok := exchangeinfo.GetBaseAssetQuantityPrecision(symbol, in.OrderType == tradeengineproto.ORDER_TYPE_MARKET)
	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_execute_perpetuals_trade.asset_quantity_precision_unknown", errParams)
	}

	// Round the price to the minimum precision allowed on the exchange.
	assetPricePrecision, ok := exchangeinfo.GetBaseAssetPricePrecision(symbol)
	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_execute_perpetuals_trade.asset_price_precision_unknown", errParams)
	}

	// Parse Order ID.
	var clientOrderID string
	switch {
	case in.OrderId != "":
		clientOrderID = in.OrderId
	}

	// Convert floats to minimum precision rounded strings.
	quantity := roundToPrecisionString(float64(in.Quantity), assetQuantityPrecision)

	// Parse limit & stop price.
	var limitPrice, stopPrice string
	switch in.OrderType {
	case tradeengineproto.ORDER_TYPE_MARKET:
		// Do nothing.
	case tradeengineproto.ORDER_TYPE_LIMIT:
		limitPrice = roundToPrecisionString(float64(in.LimitPrice), assetPricePrecision)
	case tradeengineproto.ORDER_TYPE_STOP_LIMIT:
		limitPrice = roundToPrecisionString(float64(in.LimitPrice), assetPricePrecision)
		stopPrice = roundToPrecisionString(float64(in.StopPrice), assetPricePrecision)
	case tradeengineproto.ORDER_TYPE_STOP_MARKET:
		stopPrice = roundToPrecisionString(float64(in.StopPrice), assetPricePrecision)
	case tradeengineproto.ORDER_TYPE_TAKE_PROFIT_LIMIT:
		limitPrice = roundToPrecisionString(float64(in.LimitPrice), assetPricePrecision)
		stopPrice = roundToPrecisionString(float64(in.StopPrice), assetPricePrecision)
	case tradeengineproto.ORDER_TYPE_TAKE_PROFIT_MARKET:
		stopPrice = roundToPrecisionString(float64(in.StopPrice), assetPricePrecision)
	default:
		return nil, gerrors.Unimplemented("failed_to_marshall_perpetuals_trade.unimplemented.order_type", nil)
	}

	// Parse reduce only.
	var reduceOnly string
	if in.ReduceOnly {
		reduceOnly = "true"
	}

	errParams["limit_price"] = limitPrice
	errParams["stop_price"] = stopPrice
	errParams["reduce_only"] = reduceOnly

	// Parse side.
	var side string
	switch in.TradeSide {
	case tradeengineproto.TRADE_SIDE_BUY, tradeengineproto.TRADE_SIDE_LONG:
		side = "BUY"
	case tradeengineproto.TRADE_SIDE_SELL, tradeengineproto.TRADE_SIDE_SHORT:
		side = "SELL"
	default:
		errParams["trade_side"] = side
		return nil, gerrors.BadParam("failed_to_marshall_perpetuals_trade.invalid_trade_side", errParams)
	}

	// Parse order type.
	var orderType string
	switch in.OrderType {
	case tradeengineproto.ORDER_TYPE_LIMIT:
		orderType = "LIMIT"
	case tradeengineproto.ORDER_TYPE_MARKET:
		orderType = "MARKET"
	case tradeengineproto.ORDER_TYPE_STOP_MARKET:
		orderType = "STOP_MARKET"
	case tradeengineproto.ORDER_TYPE_STOP_LIMIT:
		orderType = "STOP"
	case tradeengineproto.ORDER_TYPE_TAKE_PROFIT_LIMIT:
		orderType = "TAKE_PROFIT"
	case tradeengineproto.ORDER_TYPE_TAKE_PROFIT_MARKET:
		orderType = "TAKE_PROFIT_MARKET"
	default:
		errParams["order_type"] = orderType
		return nil, gerrors.BadParam("failed_to_marshall_perpetuals_trade.invalid_order_type", errParams)
	}

	// Parse position side.
	var positionSide string
	switch {
	// TODO: do we also allow for hedge mode at some point in the future?
	default:
		positionSide = "BOTH"
	}

	// Parse time in force.
	var timeInForce string
	if in.OrderType == tradeengineproto.ORDER_TYPE_LIMIT {
		switch in.TimeInForce {
		case tradeengineproto.TIME_IN_FORCE_GOOD_TILL_CANCELLED:
			timeInForce = "GTC"
		case tradeengineproto.TIME_IN_FORCE_FILL_OR_KILL:
			timeInForce = "FOK"
		case tradeengineproto.TIME_IN_FORCE_GOOD_TILL_CROSSING:
			timeInForce = "GTX"
		case tradeengineproto.TIME_IN_FORCE_IMMEDIATE_OR_CANCEL:
			timeInForce = "IOC"
		default:
			timeInForce = "GTC" // default value
		}
	}

	// Parse working type.
	var workingType string
	switch in.WorkingType {
	case tradeengineproto.WORKING_TYPE_CONTRACT_PRICE:
		workingType = "CONTRACT_PRICE"
	case tradeengineproto.WORKING_TYPE_MARK_PRICE:
		workingType = "MARK_PRICE"
	}

	// Marshal into dto.
	return &client.ExecutePerpetualFuturesOrderRequest{
		NewClientOrderID: clientOrderID,
		Symbol:           symbol,
		Side:             side,
		OrderType:        orderType,
		PositionSide:     positionSide,
		TimeInForce:      timeInForce,
		LimitPrice:       limitPrice,
		StopPrice:        stopPrice,
		Quantity:         quantity,
		ReduceOnly:       reduceOnly,
		ClosePosition:    strconv.FormatBool(in.ClosePosition),
		WorkingType:      workingType,
		NewOrderRespType: "ACK",
		PriceProtect:     "false",
	}, nil
}

// ProtoOrderToExecuteSpotOrderRequest ...
func ProtoOrderToExecuteSpotOrderRequest(order *tradeengineproto.Order) (*client.ExecuteSpotOrderRequest, error) {
	return nil, nil
}

func isSuccess(rsp *client.VerifyCredentialsResponse) (bool, string) {
	reasons := []string{}

	if !rsp.EnableReading {
		reasons = append(reasons, "Please enable the ability to read account")
	}

	if !rsp.EnableFutures {
		reasons = append(reasons, "Please enable futures access")
	}

	if rsp.EnableWithdrawals {
		reasons = append(reasons, "You have withdrawals enabled, please turn them off")
	}

	if !rsp.IPRestrict {
		reasons = append(reasons, "You have no ip restrictions; please consider adding them")
	}

	if !rsp.EnableSpotAndMarginTrading {
		reasons = append(reasons, "Please enable spot access")
	}

	return rsp.EnableReading && rsp.EnableFutures && rsp.EnableSpotAndMarginTrading, strings.Join(reasons, ",")
}
