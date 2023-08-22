package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetBlockHeight(client *ethclient.Client) (height uint64, err error) {
	header, err := client.HeaderByNumber(context.TODO(), nil)
	if err != nil {
		return height, err
	}
	return header.Number.Uint64(), nil
}

func GetTransactionByHash(client *ethclient.Client, txHash string) (tx *types.Transaction, err error) {
	tx, _, err = client.TransactionByHash(context.TODO(), common.HexToHash(txHash))
	if err != nil {
		return
	}
	return
}

func GetTransactionReceiptByHash(client *ethclient.Client, txHash common.Hash) (tx *types.Receipt, err error) {
	tx, err = client.TransactionReceipt(context.TODO(), txHash)
	if err != nil {
		return
	}
	return
}
