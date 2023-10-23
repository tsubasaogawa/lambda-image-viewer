package main

import (
	"fmt"
	"math"
)

func calcShutterSpeed(ssv float64) string {
	denom := int(math.Round(math.Pow(2, ssv)))
	if denom > 0 {
		return fmt.Sprintf("1/%d", denom)
	}
	return fmt.Sprintf("%d", denom)
}
