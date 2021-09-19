package handler

import (
	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/domain"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func isActorValid(actorID string) bool {
	switch actorID {
	case tradeengineproto.TradeEngineActorSatoshiSystem,
		tradeengineproto.TradeEngineActorManual,
		tradeengineproto.TradeEngineActorSatoshiCommand:
		return true
	default:
		return false
	}
}

func validateTradeParticipant(tradeParticipant *tradeengineproto.AddParticipantToTradeRequest, trade *domain.Trade) error {
	switch {
	case tradeParticipant.Risk > 50:
		return gerrors.FailedPrecondition("invalid_trade_participant.risk_too_high", nil)
	case tradeParticipant.Size < 0 && tradeParticipant.Risk < 0:
		return gerrors.BadParam("bad_param.risk_or_size_cannot_be_less_than_zero", nil)
	case tradeParticipant.Size == 0 && tradeParticipant.Risk == 0:
		return gerrors.BadParam("bad_params.risk_and_size_cannot_be_zero", nil)
	}

	return nil
}
