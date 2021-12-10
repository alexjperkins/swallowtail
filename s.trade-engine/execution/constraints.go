package execution

import "swallowtail/libraries/gerrors"

const (
	retailMaxRiskPerTradeStrategy = 10.5 // +0.5 to act as a buffer.
)

func isTradeStrategyParticipantOverRiskAppetite(accountBalance float64, totalOrderSize float64) error {
	if totalOrderSize > accountBalance*retailMaxRiskPerTradeStrategy {
		return gerrors.FailedPrecondition("trade_strategy_participant_over_risk_appetite", nil)
	}

	return nil
}
