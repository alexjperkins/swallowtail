package handler

import (
	"swallowtail/s.trade-engine/domain"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func isActorValid(actorID string) bool {
	switch actorID {
	case tradeengineproto.TradeEngineActorSatoshiSystem, tradeengineproto.TradeEngineActorManual:
		return true
	default:
		return false
	}
}

func validateTradeParticipant(tradeParticipant *tradeengineproto.AddParticipantToTradeRequest, trade *domain.Trade) error {
	return nil
}
