package mock

import (
	"math"
	"testing"
)

const (
	day = 4
)

func TestGetETDScatter(t *testing.T) {
	InitSeed(day)
	//var avg = [8]float64{0.00055, 0.0055, 0.055, 0.55, 5.5, 55, 550, 1500}
	total := GetETDNumberDay(day)
	times := GetTransactionNumberDay(day)
	exp := float64(total) / float64(times)
	result := GetETDScatter(exp)
	sum := 0.
	cost := 0.

	for j := 0; j < 1000000; j++ {
		scatter := getRandType(result[:]) - 3
		etd := GetETDNumberTransaction(scatter)
		cost += etd
	}
	cost /= 1000000

	for _, v := range result {
		sum += v
	}
	if math.Abs(sum-1) > 0.0001 {
		t.Fatal("总概率不为1", result)
	}

	sum = 0.
	if math.Abs(cost-exp) > 0.01 {
		t.Fatal("期望不符", cost, exp)
	}
	t.Log(cost, exp, result)
}

func TestGetAccountIndex(t *testing.T) {
	InitSeed(day)
	t.Log(GetAccountIndex(day))
}

func TestGetTransactionScatter(t *testing.T) {
	for i := 0; i < 50; i++ {
		t.Logf("---[day %d]---", i)
		InitSeed(int64(i))
		res := GetTransactionScatter()
		sum := 0.
		for i, v := range res {
			sum += v
			t.Log(i, v)
		}
		if math.Abs(sum-1.) > 0.12 {
			t.Fatal("总概率异常", sum)
		}
		t.Log(sum)
	}
}
