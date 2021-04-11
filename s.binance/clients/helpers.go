package clients

func calculateLimitOrderRisk(entry, sl, percentageOfAccount, accountSize float32) (float32, bool, error) {
	return 0.0, false, nil
}

func calculateDCAOrderRisk(upper, lower, sl, percentageOfAccount, accountSize float32, numDCAOrders int) ([][]float32, bool, error) {
	risks := [][]float32{}
	return risks, false, nil
}
