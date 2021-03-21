package V2

import (
	"strconv"
)

// ParseFloat converts a string to 32 bit float
func ParseFloat(number string) float32 {
	floatVar, _ := strconv.ParseFloat(number, 32)
	return float32(floatVar)
}

// ParseInt converts a string to 32 bit int
func ParseInt(number string) int {
	intVar, _ := strconv.ParseInt(number, 10, 32)
	return int(intVar)
}
