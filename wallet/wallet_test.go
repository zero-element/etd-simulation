package wallet

import (
	"fmt"
	"github.com/zero-element/etd-transaction/config"
	"github.com/zero-element/etd-transaction/rpc"
	"github.com/zero-element/go-etdereum-hdwallet"
	"log"
	"testing"
)

func TestGen(t *testing.T) {
	mnemonic := config.MF
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 800; i++ {
		rawPath := fmt.Sprintf("m/44'/60'/0'/0/%d", i)
		path := hdwallet.MustParseDerivationPath(rawPath)
		ac, err := wallet.Derive(path, true)
		if err != nil {
			t.Fatal(err)
		}
		bal, _ := rpc.BalanceAt(ac)
		t.Logf("index[%d] balance: %d", i, bal)
	}
	fmt.Print(len(wallet.Accounts()))
}

func TestChild(t *testing.T) {
	mnemonic := config.MF
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		t.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(account.Address.Hex())

	path = hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/800")
	account, err = wallet.Derive(path, false)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(account.Address.Hex())

}
