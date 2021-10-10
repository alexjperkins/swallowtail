package marshaling

import (
	"swallowtail/s.trade-engine/domain"
	tradeengineproto "swallowtail/s.trade-engine/proto"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// TradeProtoToDomain converts our proto definition of a trade to the internal domain definition.
func TradeProtoToDomain(proto *tradeengineproto.Trade) *domain.Trade {
	entries := make([]float64, 0, len(proto.Entries))
	for _, entry := range proto.Entries {
		entries = append(entries, float64(entry))
	}

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
		Entries:            entries,
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
	entries := make([]float32, 0, len(domain.Entries))
	for _, entry := range domain.Entries {
		entries = append(entries, float32(entry))
	}

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
		Entries:            entries,
		StopLoss:           float32(domain.StopLoss),
		TakeProfits:        tps,
		Status:             tradeengineproto.TRADE_STATUS((tradeengineproto.TRADE_STATUS_value[domain.Status])),
		CurrentPrice:       float32(domain.CurrentPrice),
		Created:            timestamppb.New(domain.Created),
		LastUpdated:        timestamppb.New(domain.LastUpdated),
	}
}

// TradeParticipantProtoToDomain ...
func TradeParticipantProtoToDomain(in *tradeengineproto.AddParticipantToTradeRequest, exchangeOrderID string, excutedTimestamp time.Time) *domain.TradeParticipant {
	return &domain.TradeParticipant{
		UserID:          in.UserId,
		TradeID:         in.TradeId,
		IsBot:           in.IsBot,
		Size:            float64(in.Size),
		Risk:            float64(in.Risk),
		Exchange:        in.Exchange,
		ExchangeOrderID: exchangeOrderID,
	}
}
