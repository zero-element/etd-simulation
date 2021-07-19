package mock

import (
	"etd-transaction/rpc"
	"etd-transaction/wallet"
	"github.com/ethereum/go-ethereum/accounts"
	"log"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	tsk Task
)

func TestInit(t *testing.T) {
	var tsk Task
	err := tsk.InitTask(44)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("trans per hour: %v", tsk.rTrans)
	t.Logf("etd per hour: %v", tsk.rETD)
	t.Logf("price: %d", tsk.price)
	t.Logf("ratio: %f", tsk.dRatio)
	t.Logf("account number: %d", tsk.dAccount)
}

func TestRun(t *testing.T) {
	n := time.Now()
	hour := n.Hour()
	rAct := ((59-n.Minute())*60 + (59 - n.Second())) / 10
	if tsk.rETD[hour] < 0 && hour < 23 { // 从下一小时预扣，维持最小量交易
		tsk.rETD[hour+1] += tsk.rETD[hour] - 1
		tsk.rETD[hour] = 1
	}
	if tsk.rTrans[hour] <= 0 || tsk.rETD[hour] <= 0 {
		return
	}
	var times int
	var s [8]float64
	func() {
		tsk.mu.RLock()
		defer tsk.mu.RUnlock()
		times = int(math.Round(float64(tsk.rTrans[hour]) / float64(rAct)))
		s = GetETDScatter(tsk.rETD[hour] / float64(tsk.rTrans[hour]))
	}()
	log.Printf("[%d] times: %d, ract: %d, expect: %f\nretd: %f, rtrans: %d",
		hour, times, rAct, tsk.rETD[hour]/float64(tsk.rTrans[hour]), tsk.rETD[hour], tsk.rTrans[hour])

	rChannel := make(chan float64, times)
	var wg sync.WaitGroup

	for i := 0; i < times; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var (
				from    accounts.Account
				to      accounts.Account
				balance *big.Int
				err     error
			)

			// 生成ETD
			scatter := getRandType(s[:]) - 3
			etd := GetETDNumberTransaction(scatter)
			wei := big.NewInt(int64(etd * 1e18))

			indT := rand.Intn(tsk.dAccount)
			to = wallet.GetAccountNew(indT)

			r := rand.Float64()
			if r >= tsk.dRatio {
				indF := rand.Intn(tsk.dAccount)
				for indF == indT {
					indF = rand.Intn(tsk.dAccount)
				}
				from = wallet.GetAccountNew(indF)
				balance, err = rpc.BalanceAt(from)
				if err != nil {
					log.Fatal(err.Error())
					return
				}
			}
			if r < tsk.dRatio || balance.Cmp(wei) != 1 { // 由矿工流出或者余额不足
				index := rand.Intn(wallet.GetAccountMinerNumber())
				from = wallet.GetAccountMiner(index)
			}

			t.Logf("etd: %f\nfrom: %v\nto: %v\nprice: %d", etd, from, to, tsk.price)
			rChannel <- etd
		}()
	}
	wg.Wait()

	tsk.mu.Lock()
	defer tsk.mu.Unlock()
	count := times
	sum := 0.
	for i := 0; i < times; i++ {
		if res, ok := <-rChannel; ok {
			if res == 0 {
				count -= 1
			}
			sum += res
		}
	}
	tsk.rTrans[hour] -= int64(count)
	tsk.rETD[hour] -= sum
	if tsk.rETD[hour] < 0 && hour < 23 {
		tsk.rETD[hour+1] -= tsk.rETD[hour]
		tsk.rETD[hour] = 0
	}
	t.Logf("suc times: %d", count)
	t.Logf("sum: %f", sum)
}

func TestTask_Run(t *testing.T) {
	var tsk Task
	err := tsk.InitTask(0)
	if err != nil {
		t.Fatal(err.Error())
	}
	tsk.Run()
	t.Logf("%v", tsk)
}
