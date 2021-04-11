package util

import (
	"fmt"
	"strconv"

	"github.com/monzo/terrors"
)

// FormatPriceFromString formats a string based on significant figures.
func FormatPriceFromString(price string) (string, error) {
	if price == "" {
		return "", nil
	}
	f, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "", terrors.Augment(err, "Failed to parse price string", nil)
	}
	if f > 1.0 {
		return fmt.Sprintf("%.3f", f), nil
	}

	var numZeros int
	for _, c := range price {
		if c == '0' {
			numZeros++
			continue
		}
		if c == '.' {
			continue
		}
		break
	}

	tpl := fmt.Sprintf("%%.%vf", numZeros+2)
	return fmt.Sprintf(tpl, f), nil
}

// FormatPriceAsString formats a float as a string based on significant figures.
func FormatPriceAsString(price float64) (string, error) {
	if price == 0.0 {
		return "0.00", nil
	}
	if price > 1.0 {
		return fmt.Sprintf("%.3f", price), nil
	}

	strPrice := fmt.Sprintf("%f", price)
	var numZeros int
	for _, c := range strPrice {
		if c == '0' {
			numZeros++
			continue
		}
		if c == '.' {
			continue
		}
		break
	}
	tpl := fmt.Sprintf("%%.%vf", numZeros+2)
	return fmt.Sprintf(tpl, price), nil
}
