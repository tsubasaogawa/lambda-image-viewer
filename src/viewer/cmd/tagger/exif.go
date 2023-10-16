package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func datetime2unixtime(dt string) int64 {
	if t, err := time.Parse("2006:01:02 03:04:05", dt); err != nil {
		return 0
	} else {
		return t.Unix()
	}
}

func calcShutterSpeed(ssv string) string {
	nums := strings.Split(ssv, "/")
	n, err := strconv.Atoi(nums[0])
	if err != nil {
		return ""
	}
	d, err := strconv.Atoi(nums[1])
	if err != nil {
		return ""
	}
	tv := float64(n) / float64(d)
	if tv > 0 {
		return fmt.Sprintf("1/%.0f", math.Round(math.Pow(tv, 2)))
	}
	return fmt.Sprintf("%.0f", math.Round(math.Pow(tv, 2)))
}

func calcFNumber(f string) float64 {
	return divide(f)
}

func calcFocalLength(fl string) int {
	return int(divide(fl))
}

func divide(s string) float64 {
	nums := strings.Split(s, "/")
	n, err := strconv.Atoi(nums[0])
	if err != nil {
		return 0
	}
	d, err := strconv.Atoi(nums[1])
	if err != nil {
		return 0
	}
	return math.Round(float64(n) / float64(d))
}
