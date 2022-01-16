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

func calculateTotalQuantityFromPositions(accountBalance, totalRisk float64, positions []*risk.PositionDetail) float64 {
	var f func(positions []*risk.PositionDetail) float64
	f = func(postions []*risk.PositionDetail) float64 {
		if len(postions) == 0 {
			return 0
		}

		return postions[0].RiskCoefficient + f(postions[1:])
	}

	return math.Ceil(f(positions)*accountBalance*totalRisk) / 100
}

func calculateTakeProfits(totalPositionQuantity float64, takeProfitStopPrices []float32) []*TakeProfitDetail {
	if len(takeProfitStopPrices) == 0 {
		return nil
	}

	// Calculate position to consider; we leave some amount for continuation.
	// This is `total quantity * (1 - 1/n + 1)` where `n` is the number of take profits.
	var numberOfTakeProfits = len(takeProfitStopPrices)
	positionSizeToConsider := totalPositionQuantity * (1 - (1 / float64(numberOfTakeProfits+1)))

	// Calculate quantity per take profit.
	var tpds = make([]*TakeProfitDetail, 0, len(takeProfitStopPrices))
	for _, tp := range takeProfitStopPrices {
		tpds = append(tpds, &TakeProfitDetail{
			StopPrice: float64(tp),
			Quantity:  positionSizeToConsider * 1.0 / float64(numberOfTakeProfits),
		})
	}

	return tpds
}
