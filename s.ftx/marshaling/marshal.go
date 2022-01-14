package marshaling

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.ftx/client"
	"swallowtail/s.ftx/client/auth"
	"swallowtail/s.ftx/exchangeinfo"
	ftxproto "swallowtail/s.ftx/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func VenueCredentialsProtoToFTXCredentials(credentials *tradeengineproto.VenueCredentials) *auth.Credentials {
	return &auth.Credentials{
		APIKey:     credentials.ApiKey,
		SecretKey:  credentials.SecretKey,
		Subaccount: credentials.Subaccount,
		URL:        credentials.Url,
		WSURL:      credentials.WsUrl,
	}
}

func DepositDTOToProto(deposit *client.DepositRecord) *ftxproto.DepositRecord {
	return &ftxproto.DepositRecord{
		Coin:          deposit.Coin,
		Confirmations: deposit.Confirmations,
		ConfirmedTime: timestamppb.New(deposit.ConfirmedTime),
		Fee:           float32(deposit.Fee),
		Id:            deposit.ID,
		SentTime:      timestamppb.New(deposit.SentTime),
		Size:          float32(deposit.Size),
		Status:        deposit.Status,
		Time:          timestamppb.New(deposit.Time),
		TransactionId: deposit.TXID,
	}
}

func DepositsDTOToProto(deposits []*client.DepositRecord) []*ftxproto.DepositRecord {
	protos := []*ftxproto.DepositRecord{}
	for _, d := range deposits {
		protos = append(protos, DepositDTOToProto(d))
	}

	return protos
}

func OrderProtoToDTO(order *tradeengineproto.Order) (*client.ExecuteOrderRequest, error) {
	errParams := map[string]string{
		"actor_id":   order.ActorId,
		"order_id":   order.OrderId,
		"trade_side": order.TradeSide.String(),
		"order_type": order.OrderType.String(),
		"instrument": order.Instrument,
		"asset":      order.Asset,
		"pair":       order.Pair.String(),
	}

	// Parse market.
	var market string
	switch {
	case order.Instrument == "" && order.Asset == "":
		return nil, gerrors.FailedPrecondition("missing_params.instrument_or_asset", nil)
	case order.Instrument == "":
		switch order.InstrumentType {
		case tradeengineproto.INSTRUMENT_TYPE_FUTURE_PERPETUAL:
			market = fmt.Sprintf("%s-PERP", strings.ToUpper(order.Asset))
		case tradeengineproto.INSTRUMENT_TYPE_SPOT:
			market = fmt.Sprintf("%s/%s", strings.ToUpper(order.Asset), strings.ToUpper(order.Pair.String()))
		case tradeengineproto.INSTRUMENT_TYPE_MOVE:
			return nil, gerrors.FailedPrecondition("ftx_move_contract_must_be_defined_in_the_instrument", nil)
		default:
			return nil, gerrors.Unimplemented("instrument_type.unimplemented", nil)
		}
	default:
		market = strings.ToUpper(order.Instrument)
	}

	// Parse trade side.
	var side string
	switch order.TradeSide {
	case tradeengineproto.TRADE_SIDE_BUY, tradeengineproto.TRADE_SIDE_LONG:
		side = "buy"
	case tradeengineproto.TRADE_SIDE_SELL, tradeengineproto.TRADE_SIDE_SHORT:
		side = "sell"
	default:
		return nil, gerrors.BadParam("invalid_trade_side", errParams)
	}

	// Parse order type.
	var orderType string
	switch order.OrderType {
	case tradeengineproto.ORDER_TYPE_LIMIT:
		orderType = "limit"
	case tradeengineproto.ORDER_TYPE_MARKET:
		orderType = "market"
	case tradeengineproto.ORDER_TYPE_TAKE_PROFIT_MARKET:
		orderType = "takeProfit"
	case tradeengineproto.ORDER_TYPE_STOP_MARKET:
		orderType = "stop"
	case tradeengineproto.ORDER_TYPE_TRAILING_STOP_MARKET:
		orderType = "trailingStop"
	default:
		return nil, gerrors.BadParam("invalid_order_type", map[string]string{
			"type": order.OrderType.String(),
		})
	}

	// Gather instrument data.
	exchangeInstrumentData, ok := exchangeinfo.GetInstrumentBySymbol(order.Instrument)
	if !ok {
		return nil, gerrors.FailedPrecondition("exchange_instrument_metadata_not_found", errParams)
	}

	// Parse prices.
	var (
		price        string
		triggerPrice string
		orderPrice   string
		trailValue   string
	)
	switch order.OrderType {
	case tradeengineproto.ORDER_TYPE_LIMIT:
		price = roundToPrecisionString(float64(order.LimitPrice), exchangeInstrumentData.MininumTickSize)
	case tradeengineproto.ORDER_TYPE_STOP_LIMIT, tradeengineproto.ORDER_TYPE_TAKE_PROFIT_LIMIT:
		orderPrice = roundToPrecisionString(float64(order.LimitPrice), exchangeInstrumentData.MininumTickSize)
		triggerPrice = roundToPrecisionString(float64(order.StopPrice), exchangeInstrumentData.MininumTickSize)
	case tradeengineproto.ORDER_TYPE_STOP_MARKET, tradeengineproto.ORDER_TYPE_TAKE_PROFIT_MARKET:
		triggerPrice = roundToPrecisionString(float64(order.LimitPrice), exchangeInstrumentData.MininumTickSize)
	}

	// Parse quantity.
	var quantity string
	if !order.ClosePosition {
		quantity = roundToPrecisionString(float64(order.Quantity), exchangeInstrumentData.MininumQuantity)
	}

	// Parse IOC.
	var ioc bool
	if order.TimeInForce == tradeengineproto.TIME_IN_FORCE_IMMEDIATE_OR_CANCEL {
		ioc = true
	}

	// Parse retry until filled.
	var retryUntilFilled bool
	switch order.OrderType {
	case tradeengineproto.ORDER_TYPE_STOP_MARKET, tradeengineproto.ORDER_TYPE_TAKE_PROFIT_MARKET:
		retryUntilFilled = true
	}

	// Marshal into DTO.
	return &client.ExecuteOrderRequest{
		ClientID:         order.OrderId,
		Market:           market,
		Side:             side,
		Type:             orderType,
		Price:            price,
		TriggerPrice:     triggerPrice,
		OrderPrice:       orderPrice,
		TrailValue:       trailValue,
		Size:             quantity,
		ReduceOnly:       order.ReduceOnly,
		IOC:              ioc,
		PostOnly:         order.PostOnly,
		RetryUntilFilled: retryUntilFilled,
	}, nil
}

// InstrumentsDTOToProtos ...
func InstrumentsDTOToProtos(ii []*client.Instrument) []*ftxproto.Instrument {
	protos := make([]*ftxproto.Instrument, 0, len(ii))
	for _, i := range ii {
		protos = append(protos, InstrumentDTOToProto(i))
	}

	return protos
}

// InstrumentDTOToProto ...
func InstrumentDTOToProto(i *client.Instrument) *ftxproto.Instrument {
	return &ftxproto.Instrument{
		Symbol:          i.Symbol,
		PostOnly:        i.PostOnly,
		Enabled:         i.Enabled,
		MinimumTickSize: float32(i.MininumTickSize),
		MinimumQuantity: float32(i.MininumQuantity),
		Type:            i.Type,
		Underlying:      i.Underlying,
		BaseCurrency:    i.BaseCurrency,
		QuoteCurrency:   i.QuoteCurrency,
	}
}

func ReadAccountInformationDomainToProto(in *client.ReadAccountInformationResponseResult) *ftxproto.ReadAccountInformationResponse {
	return &ftxproto.ReadAccountInformationResponse{
		BackstopProvider:             in.BackstopProvider,
		Collateral:                   float32(in.Collateral),
		FreeCollateral:               float32(in.Collateral),
		InitalMarginRequirement:      float32(in.InitialMarginRequirement),
		Leverage:                     float32(in.Leverage),
		Liquidating:                  in.Liquidating,
		MaintenanceMarginRequirement: float32(in.MaintenanceMarginRequirement),
		MakerFee:                     float32(in.MakerFee),
		MarginFraction:               float32(in.MakerFee),
		OpenMarginFraction:           float32(in.OpenMarginFraction),
		TakerFee:                     float32(in.TakerFee),
		TotalAccountValue:            float32(in.TotalAccountValue),
		TotalPositionSize:            float32(in.TotalPositionSize),
		Username:                     in.Username,
	}
}

// AccountBalancesDTOToProtos ...
func AccountBalancesDTOToProtos(in []*client.AccountBalance) []*ftxproto.AccountBalance {
	var out = make([]*ftxproto.AccountBalance, 0, len(in))
	for _, ab := range in {
		out = append(out, AccountBalancesDTOToProto(ab))
	}

	return out
}

// AccountBalancesDTOToProto ...
func AccountBalancesDTOToProto(in *client.AccountBalance) *ftxproto.AccountBalance {
	return &ftxproto.AccountBalance{
		Asset:                  in.Asset,
		QuantityAvailable:      float32(in.QuantityAvailable),
		SpotBorrow:             float32(in.SpotBorrow),
		Total:                  float32(in.Total),
		UsdValue:               float32(in.USDValue),
		AvailableWithoutBorrow: float32(in.AvailableWithoutBorrow),
	}
}

func roundToPrecisionString(f float64, minIncrement float64) string {
	v := f / minIncrement

	var p float64
	switch {
	case v < 1.0:
		p = math.Ceil(f) * minIncrement
	default:
		p = math.Floor(f) * minIncrement
	}

	// Format float & trim zeros.
	return strings.TrimRight(strconv.FormatFloat(p, 'f', 6, 64), "0")
}
