package marshaling

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	"swallowtail/s.binance/exchangeinfo"
	binanceproto "swallowtail/s.binance/proto"
)

// CredentialsProtoToDTO ...
func CredentialsProtoToDTO(in *binanceproto.Credentials) *client.Credentials {
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
func erpetualFuturesAccountBalanceDTOToProto(in *client.PerpetualFuturesAccountBalance) (*binanceproto.ReadPerpetualFuturesAccountResponse, error) {
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

// ProtoOrdersToExecutePerpetualsFutureTradeRequest ...
func ProtoOrdersToExecutePerpetualsFutureTradeRequest(ins []*binanceproto.PerpetualFuturesOrder) ([]*client.ExecutePerpetualFuturesTradeRequest, error) {
	orders := make([]*client.ExecutePerpetualFuturesTradeRequest, 0, len(ins))
	for _, in := range ins {
		order, err := ProtoOrderToExecutePerpetualsFutureTradeRequest(in)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// ProtoOrderToExecutePerpetualsFutureTradeRequest ...
func ProtoOrderToExecutePerpetualsFutureTradeRequest(in *binanceproto.PerpetualFuturesOrder) (*client.ExecutePerpetualFuturesTradeRequest, error) {
	// Round the quantity to the minimum precision allowed on the exchange.
	assetQuantityPrecision, ok := exchangeinfo.GetBaseAssetQuantityPrecision(in.Symbol)
	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_execute_perpetuals_trade.asset_quantity_precision_unknown", nil)
	}

	// Round the price to the minimum precision allowed on the exchange.
	assetPricePrecision, ok := exchangeinfo.GetBaseAssetPricePrecision(in.Symbol)
	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_execute_perpetuals_trade.asset_price_precision_unknown", nil)
	}

	// Convert floats to minimum precision rounded strings.
	quantity := roundToPrecisionString(float64(in.Quantity), assetQuantityPrecision)
	price := roundToPrecisionString(float64(in.Price), assetPricePrecision)
	stopPrice := roundToPrecisionString(float64(in.StopPrice), assetPricePrecision)

	// Parse reduce only.
	var reduceOnly string
	if !in.ClosePosition && (in.OrderType == binanceproto.BinanceOrderType_BINANCE_STOP_MARKET || in.OrderType == binanceproto.BinanceOrderType_BINANCE_TAKE_PROFIT_MARKET) {
		reduceOnly = "true"
	}

	errParams := map[string]string{}

	var side string
	switch in.Side {
	case binanceproto.BinanceTradeSide_BINANCE_BUY:
		side = "BUY"
	case binanceproto.BinanceTradeSide_BINANCE_SELL:
		side = "SELL"
	default:
		errParams["side"] = side
		return nil, gerrors.BadParam("failed_to_marshall_perpetuals_trade.invalid_order_type", errParams)
	}

	var orderType string
	switch in.OrderType {
	case binanceproto.BinanceOrderType_BINANCE_LIMIT:
		orderType = "LIMIT"
	case binanceproto.BinanceOrderType_BINANCE_MARKET:
		orderType = "MARKET"
	case binanceproto.BinanceOrderType_BINANCE_STOP_MARKET:
		orderType = "STOP_MARKET"
	case binanceproto.BinanceOrderType_BINANCE_STOP:
		orderType = "STOP"
	case binanceproto.BinanceOrderType_BINANCE_TAKE_PROFIT:
		orderType = "TAKE_PROFIT"
	case binanceproto.BinanceOrderType_BINANCE_TAKE_PROFIT_MARKET:
		orderType = "TAKE_PROFIT_MARKET"
	default:
		errParams["order_type"] = orderType
		return nil, gerrors.BadParam("failed_to_marshall_perpetuals_trade.invalid_order_type", errParams)
	}

	var positionSide string
	switch in.PositionSide {
	case binanceproto.BinancePositionSide_BINANCE_SIDE_BOTH:
		positionSide = "BOTH"
	case binanceproto.BinancePositionSide_BINANCE_SIDE_SHORT:
		positionSide = "SHORT"
	case binanceproto.BinancePositionSide_BINANCE_SIDE_LONG:
		positionSide = "LONG"
	}

	var timeInForce string
	switch in.TimeInForce {
	case binanceproto.BinanceTimeInForce_BINANCE_GTC:
		timeInForce = "GTC"
	case binanceproto.BinanceTimeInForce_BINANCE_FOK:
		timeInForce = "FOK"
	case binanceproto.BinanceTimeInForce_BINANCE_GTX:
		timeInForce = "GTX"
	case binanceproto.BinanceTimeInForce_BINANCE_IOC:
		timeInForce = "IOC"
	default:
		// Leave emtpy
	}

	var workingType string
	switch in.WorkingType {
	case binanceproto.BinanceWorkingType_BINANCE_CONTRACT_PRICE:
		workingType = "CONTRACT_PRICE"
	case binanceproto.BinanceWorkingType_BINANCE_MARK_PRICE:
		workingType = "MARKET_PRICE"
	}

	return &client.ExecutePerpetualFuturesTradeRequest{
		Symbol:           in.Symbol,
		Side:             side,
		OrderType:        orderType,
		PositionSide:     positionSide,
		TimeInForce:      timeInForce,
		Price:            price,
		StopPrice:        stopPrice,
		Quantity:         quantity,
		ReduceOnly:       reduceOnly,
		ClosePosition:    strconv.FormatBool(in.ClosePosition),
		WorkingType:      workingType,
		NewOrderRespType: "ACK",
		PriceProtect:     "false",
	}, nil
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

// NOTE: this **does** not account for large floats & can lead to overflow
func roundToPrecision(f float64, p int) float64 {
	return math.Round(f*(math.Pow10(p))) / math.Pow10(p)
}

// NOTE: this **does** not account for large floats & can lead to overflow
func roundToPrecisionString(f float64, p int) string {
	if f == 0 {
		return ""
	}

	format := fmt.Sprintf("%%.%vf", p)
	return fmt.Sprintf(format, math.Round(f*(math.Pow10(p)))/math.Pow10(p))
}
