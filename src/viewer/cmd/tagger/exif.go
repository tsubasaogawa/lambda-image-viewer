package main

import (
	"fmt"
	"math"
)

// calcShutterSpeed calculates the shutter speed based on the APEX value (ssv).
// The shutter speed is represented as the reciprocal of 2 raised to the power of the given APEX value.
// For example, an APEX value of 1 results in a shutter speed of "1/2".
func calcShutterSpeed(ssv float64) (string, error) {
	if ssv < 0 {
		return "", fmt.Errorf("ssv must be non-negative")
	}
	denom := int(math.Round(math.Pow(2, ssv)))
	if denom > 0 {
		return fmt.Sprintf("1/%d", denom), nil
	}
	return fmt.Sprintf("%d", denom), nil
}
