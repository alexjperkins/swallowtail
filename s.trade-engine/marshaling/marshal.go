package marshaling

import (
	"swallowtail/s.trade-engine/domain"
	tradeengineproto "swallowtail/s.trade-engine/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// TradeProtoToDomain converts our proto definition of a trade to the internal domain definition.
func TradeProtoToDomain(proto *tradeengineproto.Trade) *domain.Trade {
	tps := make([]float64, 0, len(proto.TakeProfits))
	for _, tp := range proto.TakeProfits {
		tps = append(tps, float64(tp))
	}

	return &domain.Trade{
		TradeID:            proto.TradeId,
		ActorID:            proto.ActorId,
		ActorType:          proto.ActorType.String(),
		HumanizedActorName: proto.HumanizedActorName,
		IdempotencyKey:     proto.IdempotencyKey,
		OrderType:          proto.OrderType.String(),
		TradeType:          proto.TradeType.String(),
		TradeSide:          proto.TradeSide.String(),
		Asset:              proto.Asset,
		Pair:               proto.Pair.String(),
		Entry:              float64(proto.Entry),
		StopLoss:           float64(proto.StopLoss),
		TakeProfits:        tps,
		CurrentPrice:       float64(proto.CurrentPrice),
		Status:             proto.Status.String(),
		Created:            proto.Created.AsTime(),
		LastUpdated:        proto.LastUpdated.AsTime(),
	}
}

// TradeDomainToProto ...
func TradeDomainToProto(domain *domain.Trade) *tradeengineproto.Trade {
	tps := make([]float32, 0, len(domain.TakeProfits))
	for _, tp := range domain.TakeProfits {
		tps = append(tps, float32(tp))
	}

	return &tradeengineproto.Trade{
		TradeId:            domain.TradeID,
		ActorId:            domain.ActorID,
		ActorType:          tradeengineproto.ACTOR_TYPE((tradeengineproto.ACTOR_TYPE_value[domain.ActorType])),
		HumanizedActorName: domain.HumanizedActorName,
		OrderType:          tradeengineproto.ORDER_TYPE((tradeengineproto.ORDER_TYPE_value[domain.OrderType])),
		TradeType:          tradeengineproto.TRADE_TYPE((tradeengineproto.TRADE_TYPE_value[domain.TradeType])),
		TradeSide:          tradeengineproto.TRADE_SIDE((tradeengineproto.TRADE_SIDE_value[domain.TradeSide])),
		Asset:              domain.Asset,
		Pair:               tradeengineproto.TRADE_PAIR((tradeengineproto.TRADE_PAIR_value[domain.Pair])),
		Entry:              float32(domain.Entry),
		StopLoss:           float32(domain.StopLoss),
		TakeProfits:        tps,
		Status:             tradeengineproto.TRADE_STATUS((tradeengineproto.TRADE_STATUS_value[domain.Status])),
		CurrentPrice:       float32(domain.CurrentPrice),
		Created:            timestamppb.New(domain.Created),
		LastUpdated:        timestamppb.New(domain.LastUpdated),
	}
}

// TradeParticipantProtoToDomain
func TradeParticipantProtoToDomain(in *tradeengineproto.AddParticipantToTradeRequest) *domain.TradeParticipent {
	return &domain.TradeParticipent{}
}
