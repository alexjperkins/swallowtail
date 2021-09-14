package risk

import (
	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// CalculateNotionalSizeFromPositionAndRisk ...
func CalculateNotionalSizeFromPositionAndRisk(entry, stopLoss, risk, accountSize float64, side *tradeengineproto.TRADE_SIDE) (float64, error) {
	if entry == stopLoss {
		return 0, gerrors.FailedPrecondition("failed_to_calculate_notional_size_from_risk.entry_cannot_equal_stop_loss", nil)
	}

	maxRiskToLose := risk * accountSize * 0.01
	lossPerContract := entry - stopLoss

	notionalSize := maxRiskToLose / lossPerContract

	switch side {
	case tradeengineproto.TRADE_SIDE_LONG.Enum(), tradeengineproto.TRADE_SIDE_BUY.Enum():
		return notionalSize, nil
	case tradeengineproto.TRADE_SIDE_SHORT.Enum(), tradeengineproto.TRADE_SIDE_SELL.Enum():
		return -1 * notionalSize, nil
	default:
		return 0, gerrors.Unimplemented("failed_to_calculate_notional_size_from_risk.trade_side_unimplemented", map[string]string{
			"side": side.String(),
		})
	}
}
