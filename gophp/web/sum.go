package web

import (
	"math"
)

// x对y求余
func Mod(x, y int) int {
	return int(math.Mod(float64(x), float64(y)))
}
