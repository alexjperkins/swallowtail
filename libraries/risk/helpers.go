package risk

func summedLinspace(howMany int) []float64 {
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

	normalizationCoeff := 1.0 / t

	for i := range vs {
		vs[i] *= normalizationCoeff
	}

	return vs
}

func abs(v float64) float64 {
	if v < 0 {
		return -1 * v
	}

	return v
}
