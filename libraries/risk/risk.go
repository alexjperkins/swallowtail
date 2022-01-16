package risk

import (
	"math"
	"sort"
	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// PositionDetail ...
type PositionDetail struct {
	// RiskCoefficient is defined as the number of contracts one needs to buy or sell per unit of risk.
	RiskCoefficient float64
	Price           float64
}

// CalculateRiskCoefficient calculates the risk coefficient given an entry & and a stop loss.
func CalculateRiskCoefficient(entry, stopLoss float64) float64 {
	return 0.01 / abs(entry-stopLoss)
}

// CalculatePositionsByRisk returns an array of risk percentages to  place based on both the entry, position & total risk.
func CalculatePositionsByRisk(
	entries []float64,
	stopLoss float64,
	howMany int,
	side tradeengineproto.TRADE_SIDE,
	strategy tradeengineproto.DCA_EXECUTION_STRATEGY,
) ([]*PositionDetail, error) {
	risks := make([]*PositionDetail, 0, len(entries))

	// Calculate position sizes.
	positions := make([]float64, 0, len(entries))
	switch len(entries) {
	case 0:
		return nil, gerrors.FailedPrecondition("failed_to_calculate_notional_size_from_risk.missing_entries", nil)
	case 1:
		r, err := calculateContractsPerUnitRisk(entries[0], stopLoss, 1, side)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_calculate_notional_size_from_risk.risk_calc", nil)
		}

		return []*PositionDetail{
			{
				RiskCoefficient: r,
				Price:           entries[0],
			},
		}, nil
	default:
		var sideCoeff float64
		switch side {
		case tradeengineproto.TRADE_SIDE_BUY, tradeengineproto.TRADE_SIDE_LONG:
			sideCoeff = 1.0
		case tradeengineproto.TRADE_SIDE_SELL, tradeengineproto.TRADE_SIDE_SHORT:
			sideCoeff = -1.0
		default:
			sideCoeff = 1.0
		}

		// (Last entry - first entry) * side coeff * 1 / (num of bids - 1)
		positionIncrement := float64(math.Abs(entries[len(entries)-1]-entries[0])) * 1.0 / (float64(howMany) - 1.0) * sideCoeff
		for i := 0; i < howMany; i++ {
			positions = append(positions, entries[0]+(float64(i)*positionIncrement))
		}
	}

	// Calculate the given risk space we shall use for our risk coefficients
	var riskSpace []float64
	switch strategy {
	case tradeengineproto.DCA_EXECUTION_STRATEGY_CONSTANT:
		for i := 0; i < howMany; i++ {
			riskSpace = append(riskSpace, 1/float64(howMany))
		}
	case tradeengineproto.DCA_EXECUTION_STRATEGY_LINEAR:
		riskSpace = summedLinspace(howMany)
	case tradeengineproto.DCA_EXECUTION_STRATEGY_EXPONENTIAL:
		return nil, gerrors.Unimplemented("failed_to_calculate_notional_size_from_risk.exponential_dca_strategy_unimplemented", nil)
	default:
	}

	// Reverse the risk space.
	sort.Slice(riskSpace, func(i, j int) bool {
		return riskSpace[i] > riskSpace[j]
	})

	// Calculate risk array.
	for i, position := range positions {
		coeff := riskSpace[i]
		risk, err := calculateContractsPerUnitRisk(position, stopLoss, coeff, side)
		if err != nil {
			return nil, err
		}

		risks = append(risks, &PositionDetail{
			RiskCoefficient: risk,
			Price:           position,
		})
	}

	return risks, nil
}

func calculateContractsPerUnitRisk(entry float64, stopLoss, coeff float64, side tradeengineproto.TRADE_SIDE) (float64, error) {
	if entry == stopLoss {
		return 0, gerrors.FailedPrecondition("failed_to_calculate_contracts_per_unit_risk.entry_cannot_equal_stop_loss", nil)
	}

	contractsPerUnitOfRisk := coeff / (entry - stopLoss)

	switch side.String() {
	case "LONG", "BUY":
		return contractsPerUnitOfRisk, nil
	case "SHORT", "SELL":
		return -1 * contractsPerUnitOfRisk, nil
	default:
		return 0, gerrors.Unimplemented("failed_to_calculate_contracts_per_unit_risk.trade_side_unimplemented", map[string]string{
			"side": side.String(),
		})
	}
}
