package rpc

import (
	"context"
	"etd-transaction/config"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
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
	log.Printf("ID: %d", ChainID)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func BalanceAt(account accounts.Account) (*big.Int, error) {
	res, err := c.BalanceAt(context.Background(), account.Address, nil)
	if err != nil {
		log.Print(err.Error(), account)
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
