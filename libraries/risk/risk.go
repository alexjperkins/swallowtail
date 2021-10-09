package risk

import (
	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// CalculateNotionalSizesFromPositionAndRisk returns an array of risk percentages to  place based on both the entry, position & total risk
func CalculateNotionalSizesFromPositionAndRisk(entries []float64, stopLoss, totalRisk float64, side *tradeengineproto.TRADE_SIDE) ([]float64, error) {
	risks := make([]float64, 0, len(entries))
	for _, entry := range entries {
		risk, err := calculateRisk(entry, stopLoss, totalRisk, side)
		if err != nil {
			return nil, err
		}

		risks = append(risks, risk)
	}

	return risks, nil
}

func calculateRisk(entry float64, stopLoss, risk float64, side *tradeengineproto.TRADE_SIDE) (float64, error) {
	if entry == stopLoss {
		return 0, gerrors.FailedPrecondition("failed_to_calculate_notional_size_from_risk.entry_cannot_equal_stop_loss", nil)
	}

	perc := risk * 0.01
	lossPerContract := entry - stopLoss

	notionalSize := perc / lossPerContract

	switch side.String() {
	case "LONG", "BUY":
		return notionalSize, nil
	case "SHORT", "SELL":
		return -1 * notionalSize, nil
	default:
		return 0, gerrors.Unimplemented("failed_to_calculate_notional_size_from_risk.trade_side_unimplemented", map[string]string{
			"side": side.String(),
		})
	}
}
