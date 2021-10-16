package marshaling

import (
	"fmt"
	"math"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.ftx/client"
	ftxproto "swallowtail/s.ftx/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func OrdersProtoToDTO(orders []*ftxproto.FTXOrder) ([]*client.ExecuteOrderRequest, error) {
	protos := make([]*client.ExecuteOrderRequest, 0, len(orders))
	for _, o := range orders {
		proto, err := OrderProtoToDTO(o)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_marshal_orders_to_dto", nil)
		}

		protos = append(protos, proto)
	}

	return protos, nil
}

func OrderProtoToDTO(order *ftxproto.FTXOrder) (*client.ExecuteOrderRequest, error) {
	var side string
	switch order.Side {
	case ftxproto.FTX_SIDE_FTX_SIDE_BUY:
		side = "buy"
	case ftxproto.FTX_SIDE_FTX_SIDE_SELL:
		side = "sell"
	default:
		return nil, gerrors.FailedPrecondition("unrecognized_side", map[string]string{
			"side": order.Side.String(),
		})
	}

	var tradeType string
	switch order.Type {
	case ftxproto.FTX_TRADE_TYPE_FTX_TRADE_TYPE_LIMIT:
		tradeType = "limit"
	case ftxproto.FTX_TRADE_TYPE_FTX_TRADE_TYPE_MARKET:
		tradeType = "markte"
	case ftxproto.FTX_TRADE_TYPE_FTX_TRADE_TYPE_TAKE_PROFIT:
		tradeType = "takeProfit"
	case ftxproto.FTX_TRADE_TYPE_FTX_TRADE_TYPE_STOP:
		tradeType = "stop"
	case ftxproto.FTX_TRADE_TYPE_FTX_TRADE_TYPE_TRALING_STOP:
		tradeType = "trailingStop"
	default:
		return nil, gerrors.FailedPrecondition("unrecognized_type", map[string]string{
			"type": order.Type.String(),
		})
	}

	var (
		price        string
		quantity     string
		triggerPrice string
		orderPrice   string
		trailValue   string
	)

	if order.Price > 0 {

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
