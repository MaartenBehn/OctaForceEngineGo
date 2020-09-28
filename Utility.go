package OctaForceEngine

import (
	"strconv"
)

func ParseFloat(number string) float32 {
	floatVar, _ := strconv.ParseFloat(number, 32)
	return float32(floatVar)
}
func ParseInt(number string) int {
	intVar, _ := strconv.ParseInt(number, 10, 32)
	return int(intVar)
}
