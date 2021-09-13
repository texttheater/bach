// Package varhypot provides a variadic version of the mathematical hypot
// function, accepting an arbitrary number of arguments.
package varhypot

import (
	"math"
)

// Hypot returns Sqrt(x1*x1, x2*x2...), taking care to avoid unnecessary
// overflow and underflow.
func Hypot(arguments ...float64) float64 {
	// Implementation note: this is line-by-line translated from
	// https://github.com/v8/v8/blob/fdc9fade975d9114ead769e7aeb4fd9a490f80f8/src/builtins/math.tq#L403

	// Find max and deal with edge cases.
	length := len(arguments)
	if length == 0 {
		return 0
	}
	absValues := make([]float64, length)
	oneArgIsNaN := false
	max := 0.0
	for i, value := range arguments {
		if math.IsNaN(value) {
			oneArgIsNaN = true
		} else {
			absValue := math.Abs(value)
			absValues[i] = absValue
			if absValue > max {
				max = absValue
			}
		}
	}
	if math.IsInf(max, 1) {
		return math.Inf(1)
	} else if oneArgIsNaN {
		return math.NaN()
	} else if max == 0 {
		return 0
	}

	// Kahan summation to avoid rounding errors.
	// Normalize the numbers to the largest one to avoid rounding errors.
	sum := 0.0
	compensation := 0.0
	for _, value := range absValues {
		n := value / max
		summand := n*n - compensation
		preliminary := sum + summand
		compensation = (preliminary - sum) - summand
		sum = preliminary
	}
	return math.Sqrt(sum) * max
}
