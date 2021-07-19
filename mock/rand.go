package mock

import (
	"math"
	"math/rand"
)

// 根据天数初始化种子 day从1开始
func InitSeed(day int64) {
	rand.Seed(day)
}

func getFloat64Rand(l, r float64) float64 {
	return l + rand.Float64()*(r-l)
}

func normalFloat64(x int64, center int64, sigma int64) float64 {
	randomNormal := 1 / (math.Sqrt(2*math.Pi) * float64(sigma)) * math.Pow(math.E, -math.Pow(float64(x-center), 2)/(2*math.Pow(float64(sigma), 2)))
	return randomNormal
}

func getRandType(s []float64) int {
	sum := 0.
	for _, v := range s {
		sum += v
	}

	r := getFloat64Rand(0, sum)

	sum = 0
	for i, v := range s {
		sum += v
		if sum > r {
			return i
		}
	}
	return len(s) - 1
}
