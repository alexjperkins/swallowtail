package risk

import (
	"math"
	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// RiskCalculatedPosition
type RiskCalculatedPosition struct {
	Risk  float64
	Price float64
}

// CalculatePositionsByRisk returns an array of risk percentages to  place based on both the entry, position & total risk
func CalculatePositionsByRisk(
	entries []float64,
	stopLoss, totalRisk float64,
	howMany int,
	side *tradeengineproto.TRADE_SIDE,
	strategy *tradeengineproto.DCA_STRATEGY,
) ([]*RiskCalculatedPosition, error) {

	risks := make([]*RiskCalculatedPosition, 0, len(entries))

	// Calculate position sizes
	positions := make([]float64, 0, len(entries))
	switch len(entries) {
	case 0:
		return nil, gerrors.FailedPrecondition("failed_to_calculate_notional_size_from_risk.missing_entries", nil)
	case 1:
		// If we only have one entry; then the risk size is 100% at that price.
		return []*RiskCalculatedPosition{
			{
				Risk:  1.0,
				Price: entries[0],
			},
		}, nil

	default:
		positionIncrement := float64(math.Abs(entries[len(entries)-1]-entries[0])) * 1.0 / float64(howMany)
		for i := 0; i < howMany; i++ {
			positions = append(positions, entries[0]+(float64(i)*positionIncrement))
		}
	}

	// Calculate the given risk space we shall use for our risk coefficients
	var riskSpace []float64
	switch strategy {
	case tradeengineproto.DCA_STRATEGY_CONSTANT.Enum():
		for i := 0; i < howMany; i++ {
			riskSpace = append(riskSpace, 1/float64(howMany))
		}
	case tradeengineproto.DCA_STRATEGY_LINEAR.Enum():
		riskSpace = summedLinspace(howMany, 1)
	case tradeengineproto.DCA_STRATEGY_EXPONENTIAL.Enum():
		return nil, gerrors.Unimplemented("failed_to_calculate_notional_size_from_risk.exponential_dca_strategy_unimplemented", nil)
	}

	// Calculate risk array.
	for i, entry := range entries {
		position, coeff := positions[i], riskSpace[i]
		risk, err := calculateRisk(entry, stopLoss, totalRisk, coeff, side)
		if err != nil {
			return nil, err
		}

		risks = append(risks, &RiskCalculatedPosition{
			Risk:  risk,
			Price: position,
		})
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
		t          float64
		start, end = 1.0, 3.0
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
