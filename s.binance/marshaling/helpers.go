package marshaling

import (
	"fmt"
	"math"
)

// NOTE: this **does** not account for large floats & can lead to overflow
func roundToPrecision(f float64, p int) float64 {
	return math.Round(f*(math.Pow10(p))) / math.Pow10(p)
}

// NOTE: this **does** not account for large floats & can lead to overflow
func roundToPrecisionString(f float64, p int) string {
	if f == 0 {
		return ""
	}

	format := fmt.Sprintf("%%.%vf", p)
	return fmt.Sprintf(format, math.Round(f*(math.Pow10(p)))/math.Pow10(p))
}
