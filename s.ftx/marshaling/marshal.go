package marshaling

import (
	"fmt"
	"math"
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
	errParams := map[string]string{}

	// Parse trade side.
	var side string
	switch order.TradeSide {
	case tradeengineproto.TRADE_SIDE_BUY, tradeengineproto.TRADE_SIDE_LONG:
		side = "buy"
	case tradeengineproto.TRADE_SIDE_SELL, tradeengineproto.TRADE_SIDE_SHORT:
		side = "sell"
	default:
		return nil, gerrors.BadParam("invalid_trade_side", map[string]string{
			"side": order.TradeSide.String(),
		})
	}

	// Parse trade type.
	var tradeType string
	switch order.OrderType {
	case tradeengineproto.ORDER_TYPE_LIMIT:
		tradeType = "limit"
	case tradeengineproto.ORDER_TYPE_MARKET:
		tradeType = "market"
	case tradeengineproto.ORDER_TYPE_TAKE_PROFIT_MARKET:
		tradeType = "takeProfit"
	case tradeengineproto.ORDER_TYPE_STOP_MARKET:
		tradeType = "stop"
	case tradeengineproto.ORDER_TYPE_TRAILING_STOP_MARKET:
		tradeType = "trailingStop"
	default:
		return nil, gerrors.BadParam("invalid_order_type", map[string]string{
			"type": order.OrderType.String(),
		})
	}

	exchangeInstrumentData, ok := exchangeinfo.GetInstrumentBySymbol(order.Instrument)
	if !ok {
		return nil, gerrors.FailedPrecondition("exchange_instrument_metadata_not_found", errParams)
	}

	// Parse order type.
	var (
		price        string
		quantity     string
		triggerPrice string
		orderPrice   string
		trailValue   string
	)
	switch order.OrderType {
	case tradeengineproto.ORDER_TYPE_LIMIT:
		price = roundToPrecisionString(float64(order.LimitPrice))
	}

	return &client.ExecuteOrderRequest{
		Side:              side,
		Type:              tradeType,
		Price:             price,
		TriggerPrice:      triggerPrice,
		OrderPrice:        orderPrice,
		TrailValue:        trailValue,
		Quantity:          quantity,
		ReduceOnly:        order.ReduceOnly,
		IOC:               order.Ioc,
		PostOnly:          order.PostOnly,
		RejectOnPriceBand: order.RejectOnPriceBand,
		RetryUntilFilled:  order.RetryUntilFilled,
	}, nil
}

// InstrumentsDTOToProtos ...
func InstrumentsDTOToProtos(is []*client.Instrument) []*ftxproto.Instrument {
	protos := make([]*ftxproto.Instrument, 0, len(is))
	for _, i := range is {
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

// NOTE: this **does** not account for large floats & can lead to overflow
// TODO: move to own library.
func roundToPrecisionString(f float64, p int) string {
	if f == 0 {
		return ""
	}

	format := fmt.Sprintf("%%.%vf", p)
	return fmt.Sprintf(format, math.Round(f*(math.Pow10(p)))/math.Pow10(p))
}
