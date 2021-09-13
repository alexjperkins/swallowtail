package handler

import (
	"swallowtail/libraries/emojis"
	satoshiproto "swallowtail/s.satoshi/proto"
)

func validateActor(actorID string) bool {
	switch actorID {
	case satoshiproto.SatoshiActorCron, satoshiproto.SatoshiActorManual, satoshiproto.SatoshiActorSatoshiSystem:
		return true
	default:
		return false
	}
}

func isValidTradeParticipantReaction(s string) bool {
	return emojis.SatoshiRiskEmoji(s).AsRiskPercentage() != 0
}
