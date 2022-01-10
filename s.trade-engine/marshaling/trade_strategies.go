package marshaling

import (
	"swallowtail/s.trade-engine/domain"
	tradeengineproto "swallowtail/s.trade-engine/proto"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// TradeProtoToDomain converts our proto definition of a trade to the internal domain definition.
func TradeStrategyProtoToDomain(proto *tradeengineproto.TradeStrategy) *domain.TradeStrategy {
	entries := make([]float64, 0, len(proto.Entries))
	for _, entry := range proto.Entries {
		entries = append(entries, float64(entry))
	}

	tps := make([]float64, 0, len(proto.TakeProfits))
	for _, tp := range proto.TakeProfits {
		tps = append(tps, float64(tp))
	}

	tradeableVenues := make([]string, 0, len(proto.TradeableVenues))
	for _, tv := range proto.TradeableVenues {
		tradeableVenues = append(tradeableVenues, tv.String())
	}

	return &domain.TradeStrategy{
		TradeStrategyID:    proto.TradeStrategyId,
		ActorID:            proto.ActorId,
		ActorType:          proto.ActorType.String(),
		HumanizedActorName: proto.HumanizedActorName,
		IdempotencyKey:     proto.IdempotencyKey,
		ExecutionStrategy:  proto.ExecutionStrategy.String(),
		InstrumentType:     proto.InstrumentType.String(),
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
		TradeableVenues:    tradeableVenues,
	}
}

// TradeStrategyDomainToProto ...
func TradeStrategyDomainToProto(domain *domain.TradeStrategy) *tradeengineproto.TradeStrategy {
	entries := make([]float32, 0, len(domain.Entries))
	for _, entry := range domain.Entries {
		entries = append(entries, float32(entry))
	}

	tps := make([]float32, 0, len(domain.TakeProfits))
	for _, tp := range domain.TakeProfits {
		tps = append(tps, float32(tp))
	}

	return &tradeengineproto.TradeStrategy{
		TradeStrategyId:    domain.TradeStrategyID,
		ActorId:            domain.ActorID,
		ActorType:          tradeengineproto.ACTOR_TYPE((tradeengineproto.ACTOR_TYPE_value[domain.ActorType])),
		HumanizedActorName: domain.HumanizedActorName,
		ExecutionStrategy:  tradeengineproto.EXECUTION_STRATEGY((tradeengineproto.EXECUTION_STRATEGY_value[domain.ExecutionStrategy])),
		InstrumentType:     tradeengineproto.INSTRUMENT_TYPE((tradeengineproto.INSTRUMENT_TYPE_value[domain.InstrumentType])),
		TradeSide:          tradeengineproto.TRADE_SIDE((tradeengineproto.TRADE_SIDE_value[domain.TradeSide])),
		Asset:              domain.Asset,
		Pair:               tradeengineproto.TRADE_PAIR((tradeengineproto.TRADE_PAIR_value[domain.Pair])),
		Entries:            entries,
		StopLoss:           float32(domain.StopLoss),
		TakeProfits:        tps,
		Status:             tradeengineproto.TRADE_STRATEGY_STATUS((tradeengineproto.TRADE_STRATEGY_STATUS_value[domain.Status])),
		CurrentPrice:       float32(domain.CurrentPrice),
		Created:            timestamppb.New(domain.Created),
		LastUpdated:        timestamppb.New(domain.LastUpdated),
	}
}

// TradeParticipantProtoToDomain ...
func TradeParticipantProtoToDomain(in *tradeengineproto.ExecuteTradeStrategyForParticipantRequest, exchangeOrderIDs []string, excutedTimestamp time.Time) *domain.TradeStrategyParticipant {
	return &domain.TradeStrategyParticipant{
		UserID:           in.UserId,
		TradeStrategyID:  in.TradeStrategyId,
		IsBot:            in.IsBot,
		Size:             float64(in.Size),
		Risk:             float64(in.Risk),
		Venue:            in.Venue.String(),
		ExchangeOrderIDs: exchangeOrderIDs,
	}
}
