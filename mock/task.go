package mock

import (
	"etd-transaction/rpc"
	"etd-transaction/wallet"
	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	log "github.com/sirupsen/logrus"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	Day      int64       //第k日
	ETD      [24]float64 //当日剩余交易额
	Trans    [24]int64   //当日剩余交易笔数
	rETD     [24]float64 //当日剩余交易额
	rTrans   [24]int64   //当日剩余交易笔数
	dRatio   float64     //当日由矿工转出比例
	dAccount int         //当日非矿机钱包总数
	price    *big.Int    //当日手续费价格
	mu       sync.RWMutex
}

func (t *Task) InitTask(day int64) error {
	var err error

	t.Day = day
	InitSeed(day)
	etd := float64(GetETDNumberDay(day))
	es := GetTransactionScatter()
	tran := GetTransactionNumberDay(day)
	for i := 0; i < 24; i++ {
		t.rETD[i] = etd * es[i]
		t.rTrans[i] = int64(float64(tran) * es[i])
	}
	t.dAccount = int(GetAccountNumberDay(day))
	t.dRatio = GetTransactionRatio()
	t.price, err = rpc.SuggestGasPrice()
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) Run() {
	n := time.Now()
	hour := n.Hour()
	rAct := ((59-n.Minute())*60+(59-n.Second()))/10 + 1
	if t.rTrans[hour] <= 0 || t.rETD[hour] <= 0 {
		return
	}
	var times int
	var s [8]float64
	func() {
		t.mu.RLock()
		defer t.mu.RUnlock()
		times = int(math.Round(float64(t.rTrans[hour]) / float64(rAct)))
		s = GetETDScatter(t.rETD[hour] / float64(t.rTrans[hour]))
	}()
	log.Infof("[%d] times: %d, ract: %d, expect: %f\nretd: %f, rtrans: %d",
		hour, times, rAct, t.rETD[hour]/float64(t.rTrans[hour]), t.rETD[hour], t.rTrans[hour])

	sc := make(chan float64, times)
	var wg sync.WaitGroup

	for i := 0; i < times; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var (
				from    accounts.Account
				to      accounts.Account
				w       *hdwallet.Wallet
				balance *big.Int
				err     error
			)

			// 生成ETD
			scatter := getRandType(s[:]) - 3
			etd := GetETDNumberTransaction(scatter)
			wei := big.NewInt(int64(etd * 1e18))

			indT := rand.Intn(t.dAccount)
			to = wallet.GetAccountNew(indT)

			r := rand.Float64()
			if r >= t.dRatio {
				w = wallet.WNew
				indF := rand.Intn(t.dAccount)
				for indF == indT {
					indF = rand.Intn(t.dAccount)
				}
				from = wallet.GetAccountNew(indF)
				balance, err = rpc.BalanceAt(from)
				if err != nil {
					log.Error(err.Error(), from)
					return
				}
			}
			if r < t.dRatio || balance.Cmp(wei) != 1 { // 由矿工流出或者余额不足
				w = wallet.WMiner
				index := rand.Intn(wallet.GetAccountMinerNumber())
				from = wallet.GetAccountMiner(index)
			}

			if err = wallet.SendTransaction(etd, w, from, to, t.price); err != nil {
				log.Error(err.Error())
			} else {
				sc <- etd
			}
		}()
	}
	wg.Wait()
	close(sc)

	t.mu.Lock()
	defer t.mu.Unlock()
	count := 0
	sum := 0.
	for {
		if res, ok := <-sc; ok {
			count += 1
			sum += res
		} else {
			break
		}
	}
	t.rTrans[hour] -= int64(count)
	t.rETD[hour] -= sum
	t.Trans[hour] += int64(count)
	t.ETD[hour] += sum
	log.Infof("suc times: %d", count)
	log.Infof("sum: %f", sum)

	if hour == 23 {
		return
	}
	if t.rETD[hour] <= 0 { // 从下一小时预扣，维持最小量交易
		t.rETD[hour+1] += t.rETD[hour] - 1
		t.rETD[hour] = 1
	}
	if t.rTrans[hour] <= 0 {
		t.rETD[hour+1] += t.rETD[hour]
		t.rETD[hour] = 0
	}
}
