package mock

import (
	"math"
	"math/rand"
)

// 某日总交易笔数
func GetTransactionNumberDay(day int64) int64 {
	const begin int64 = 43200
	const end int64 = 86400

	eps := rand.Int63n(100) - 50

	if day >= 40 {
		return end + eps
	}

	base := begin + (end-begin)/39*day
	return base + eps
}

// 某日全部交易中由矿机流出比例
// 1/3 - 1/2
func GetTransactionRatio() float64 {
	return getFloat64Rand(0.33333, 0.5)
}

// 日内交易量分布
func GetTransactionScatter() [24]float64 {
	var res [24]float64

	for i := range res {
		res[(i+8)%24] = rand.Float64()*0.03 + normalFloat64(int64(i), 12, 6)*0.66
	}
	return res
}

// 某日单日交易额
func GetETDNumberDay(day int64) int64 {
	const begin int64 = 10000
	const end int64 = 60000

	eps := rand.Int63n(100) - 50

	if day >= 40 {
		return end + eps
	}

	base := begin + (end-begin)/39*day
	return base + eps
}

// 获取单次交易ETD数量级分布
// 数量级-3 -> 4
// 数量级为4时上限为2000
//
// E		P
// -3	 	0.05x
// -2		0.1x
// -1		x
// 0	 	y
// 1		0.05y
// 2		0.005y
// 3		0.0005y
// 4		0.0001y
// 0.55(x*(0.05*0.001+0.1*0.01+1*0.1)+y*(1*1+0.05*10+0.005*100+0.0005*1000))+0.15(y*0.0001*10000)=e
// 1.15x+1.0556y=1
// 0.0555775x+1.525y=e
func GetETDScatter(e float64) [8]float64 {
	var s [8]float64
	var fac = [8]float64{0.05, 0.1, 1, 1, 0.05, 0.005, 0.0005, 0.0001}

	const (
		facA float64 = 1.15
		facB float64 = 1.0556
		facC float64 = 0.0555775
		facD float64 = 1.525
	)

	D := facA*facD - facB*facC
	x := (1.*facD - e*facB) / D
	y := (e*facA - 1.*facC) / D

	for i := 0; i < 3; i++ {
		s[i] = x * fac[i]
	}
	for i := 3; i < 8; i++ {
		s[i] = y * fac[i]
	}
	return s
}

func GetETDNumberTransaction(scale int) float64 {
	var m float64

	if scale == 4 {
		m = getFloat64Rand(0.1, 0.2)
	} else {
		m = getFloat64Rand(0.1, 1)
	}
	return m * math.Pow10(scale)
}

// 获取某日待交易非矿机账户数量
func GetAccountNumberDay(day int64) int64 {
	const begin int64 = 50
	const end int64 = 800

	eps := rand.Int63n(10) - 5

	if day >= 40 {
		return end + eps
	}

	base := begin + (end-begin)/39*day
	return base + eps
}

// 获取某日待交易非矿机账户下标序列
func GetAccountIndex(day int64) []int {
	set := make(map[int]bool)
	l := GetAccountNumberDay(day)
	list := make([]int, l)
	for i := range list {
		for {
			r := rand.Intn(800)
			if set[r] == false {
				set[r] = true
				list[i] = r
				break
			}
		}
	}
	return list
}
