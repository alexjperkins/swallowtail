package execution

import "swallowtail/libraries/gerrors"

const (
	retailMaxRiskPerTradeStrategy = 10.5 // +0.5 to act as a buffer.
	retailMinVenueMargainInUSDT   = 100
)

func isTradeStrategyParticipantOverRiskAppetite(accountBalance float64, totalOrderSize float64) error {
	if totalOrderSize > accountBalance*retailMaxRiskPerTradeStrategy {
		return gerrors.FailedPrecondition("trade_strategy_participant_over_risk_appetite", nil)
	}

	return nil
}

func isEnoughAvailableVenueMargain(accountBalanceInUSDT float64) error {
	if accountBalanceInUSDT < retailMinVenueMargainInUSDT {
		return gerrors.FailedPrecondition("venue_balance_too_small", nil)
	}

	return nil
}
