package execution

import (
	"math"
	"swallowtail/libraries/risk"
)

const (
	DCANumberOfBuysLowerBound = 5
	DCANumberOfBuysUpperBound = 7
)

type TakeProfitDetail struct {
	StopPrice float64
	Quantity  float64
}

func calculateNumberOfDCABuys(accountBalance float64) int {
	if accountBalance > 1000 {
		return DCANumberOfBuysUpperBound
	}

	return DCANumberOfBuysLowerBound
}

func calculateTotalQuantityFromPositions(accountBalance, totalRisk float64, positions []*risk.Position) float64 {
	var f func(positions []*risk.Position) float64
	f = func(postions []*risk.Position) float64 {
		if len(postions) == 0 {
			return 0
		}

		return postions[0].RiskCoefficient + f(postions[1:])
	}

	return math.Ceil(f(positions) * accountBalance * totalRisk)
}

func calculateTakeProfits(totalPositionQuantity float64, takeProfitStopPrices []float32) []*TakeProfitDetail {
	if len(takeProfitStopPrices) == 0 {
		return nil
	}

	// Calculate position to consider; we leave some amount for continuation.
	positionSizeToConsider := totalPositionQuantity * float64(len(takeProfitStopPrices))

	var tpds = make([]*TakeProfitDetail, 0, len(takeProfitStopPrices))
	for _, tp := range takeProfitStopPrices {
		tpds = append(tpds, &TakeProfitDetail{
			StopPrice: float64(tp),
			Quantity:  positionSizeToConsider * float64(1.0/len(takeProfitStopPrices)),
		})
	}

	return tpds
}
