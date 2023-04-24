package max

import (
	"math"
)

// Max returns the maximum value found in x.
func Max(x []int) int {
	max := math.MinInt
	for _, value := range x {
		if value > max {
			max = value
		}
	}
	return max
}
