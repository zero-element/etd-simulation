package wallet

import (
	"etd-transaction/config"
	"etd-transaction/rpc"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"log"
	"math/big"
)

var (
	adMiner []accounts.Account
	adNew   []accounts.Account

	WMiner *hdwallet.Wallet
	WNew   *hdwallet.Wallet
)

func init() {
	var err error
	mnemonicFrom := config.MF
	WMiner, err = hdwallet.NewFromMnemonic(mnemonicFrom)
	if err != nil {
		panic(err)
	}

	mnemonicTo := config.MT
	WNew, err = hdwallet.NewFromMnemonic(mnemonicTo)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 200; i++ {
		rawPath := fmt.Sprintf("m/44'/60'/0'/0/%d", i)
		path := hdwallet.MustParseDerivationPath(rawPath)
		account, err := WMiner.Derive(path, true)
		if err != nil {
			panic(err)
		}

		balance, err := rpc.BalanceAt(account)
		if err != nil {
			panic(err)
		}
		if balance.Cmp(big.NewInt(0)) == 1 {
			adMiner = append(adMiner, account)
		}
	}
	log.Printf("Number of usable miner address: %d", len(adMiner))

	for i := 0; i < 800; i++ {
		rawPath := fmt.Sprintf("m/44'/60'/0'/0/%d", i)
		path := hdwallet.MustParseDerivationPath(rawPath)
		account, err := WNew.Derive(path, true)
		if err != nil {
			panic(err)
		}
		adNew = append(adNew, account)
	}
}

func GetAccountMiner(index int) accounts.Account {
	return adMiner[index]
}

func GetAccountNew(index int) accounts.Account {
	return adNew[index]
}

func GetAccountMinerNumber() int {
	return len(adMiner)
}

func SendTransaction(etd float64, w *hdwallet.Wallet, from, to accounts.Account, price *big.Int) error {
	if etd > 2000 || etd <= 0 {
		return fmt.Errorf("交易额度异常: %f", etd)
	}
	var flag bool
	if w == WMiner {
		flag = false
	} else {
		flag = true
	}

	value := big.NewInt(int64(etd * 1e18))
	toAddress := to.Address
	gasLimit := uint64(21000)
	var data []byte

	nonce, err := rpc.PendingNonceAt(from)
	if err != nil {
		log.Print(err.Error(), from, nonce)
		return err
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: price,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	})
	tx, err = w.SignTxEIP155(from, tx, rpc.ChainID, nil)
	if err != nil {
		log.Print(err.Error(), from)
		return err
	}
	log.Printf("etd: %f\nfrom: %v\nto: %v\nflag: %v\nnonce: %d", etd, from, to, flag, nonce)
	err = rpc.SendTransaction(tx)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}
