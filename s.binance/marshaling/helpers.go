package marshaling

import (
	"math"
	"strconv"
	"strings"
)

// TODO: this is copy from & tested in `s.ftx` - we should centralize this logic & somepoint
// or use a proper decimal library.
func roundToPrecisionString(f float64, minIncrement float64) string {
	if f <= 0.0 {
		return "0.0"
	}

	v := f / minIncrement

	var p float64
	switch {
	case v < 1.0:
		p = math.Ceil(v) * minIncrement
	default:
		p = math.Floor(v) * minIncrement
	}

	// Format float & trim zeros.
	return strings.TrimRight(strconv.FormatFloat(p, 'f', 6, 64), "0")
}
