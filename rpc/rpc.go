package rpc

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/zero-element/etd-transaction/config"
	"github.com/zero-element/go-etdereum/accounts"
	"github.com/zero-element/go-etdereum/core/types"
	"github.com/zero-element/go-etdereum/ethclient"
	"math/big"
)

var (
	c       *ethclient.Client
	ChainID *big.Int
)

func init() {
	var err error
	c, err = ethclient.Dial(config.RPCUrl)
	if err != nil {
		log.Fatal(err.Error())
	}
	ChainID, err = c.ChainID(context.Background())
	log.Infof("ID: %d", ChainID)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func BalanceAt(account accounts.Account) (*big.Int, error) {
	res, err := c.BalanceAt(context.Background(), account.Address, nil)
	if err != nil {
		log.Error(err.Error(), account)
		return nil, err
	}
	return res, nil
}

func SuggestGasPrice() (*big.Int, error) {
	price, err := c.SuggestGasPrice(context.Background())
	return price, err
}

func PendingNonceAt(ac accounts.Account) (uint64, error) {
	return c.PendingNonceAt(context.Background(), ac.Address)
}

func SendTransaction(tx *types.Transaction) error {
	return c.SendTransaction(context.Background(), tx)
}
