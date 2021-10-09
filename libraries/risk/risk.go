package risk

import (
	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// CalculateNotionalSizesFromPositionAndRisk returns an array of risk percentages to  place based on both the entry, position & total risk
func CalculateNotionalSizesFromPositionAndRisk(
	entries []float64,
	stopLoss, totalRisk float64,
	howMany int,
	side *tradeengineproto.TRADE_SIDE,
	strategy *tradeengineproto.DCA_STRATEGY,
) ([]float64, error) {

	risks := make([]float64, 0, len(entries))

	// Calculate the given space we shall use for our risk coefficients
	var space []float64
	switch strategy {
	case tradeengineproto.DCA_STRATEGY_CONSTANT.Enum():
		for i := 0; i < howMany; i++ {
			space = append(space, 1/float64(howMany))
		}
	case tradeengineproto.DCA_STRATEGY_LINEAR.Enum():
		space = summedLinspace(howMany, 1)
	case tradeengineproto.DCA_STRATEGY_EXPONENTIAL.Enum():
		return nil, gerrors.Unimplemented("failed_to_calculate_notional_size_from_risk.exponential_dca_strategy_unimplemented", nil)
	}

	// Calculate risk array.
	for i, entry := range entries {
		coeff := space[i]
		risk, err := calculateRisk(entry, stopLoss, totalRisk, coeff, side)
		if err != nil {
			return nil, err
		}

		risks = append(risks, risk)
	}

	return risks, nil
}

func calculateRisk(entry float64, stopLoss, risk, coefficent float64, side *tradeengineproto.TRADE_SIDE) (float64, error) {
	if entry == stopLoss {
		return 0, gerrors.FailedPrecondition("failed_to_calculate_notional_size_from_risk.entry_cannot_equal_stop_loss", nil)
	}

	perc := risk * 0.01 * coefficent
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

func summedLinspace(howMany int, total float64) []float64 {
	var (
		start float64 = 1.0
		end   float64 = 2.0
		t     float64
	)

	vs := make([]float64, 0, howMany)
	for i := 0; i < howMany; i++ {
		fi := float64(i)
		v := start + fi*(end-start)/float64(howMany)

		t += v
		vs = append(vs, v)
	}

	normalizationCoeff := total / t

	for i := range vs {
		vs[i] *= normalizationCoeff
	}

	return vs
}
