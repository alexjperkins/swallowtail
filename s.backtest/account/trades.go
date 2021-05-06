package account

import "swallowtail/s.backtest/domain"

type TradeHistory []*domain.Trade

func (t TradeHistory) StrikeRate() float64 {
	var won int
	for _, trade := range t {
		if trade.PNL > 0 {
			won++
			continue
		}
		won--
	}
	return float64(won / len(t))
}

func (t TradeHistory) FundingFees() float64 {
	var total float64
	for _, trade := range t {
		total += trade.FundingFees
	}
	return total
}

func (t TradeHistory) Fees() float64 {
	var total float64
	for _, trade := range t {
		total += trade.Fees
	}
	return total
}

func (t TradeHistory) TotalFees() float64 {
	return t.Fees() + t.FundingFees()
}

func (t TradeHistory) TotalRealizedPNL() float64 {
	var total float64
	for _, trade := range t {
		if trade.Status != domain.StatusComplete {
			continue
		}
		total += trade.PNL
	}
	return total
}
