package handlers

import (
	contracts "eth-transaction-parser/internal/transferNFT"
	"eth-transaction-parser/pkg/eth"
	"eth-transaction-parser/pkg/log"
	types2 "eth-transaction-parser/pkg/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
	"sync"
)

type (
	LogTransfer struct {
		From    common.Address
		To      common.Address
		TokenId *big.Int
	}

	LogTransferBatch struct {
		Operator common.Address
		From     common.Address
		To       common.Address
		Ids      []*big.Int
		Values   []*big.Int
	}

	LogTransferSingle struct {
		Operator common.Address
		From     common.Address
		To       common.Address
		Id       *big.Int
		Value    *big.Int
	}
)

func ParseTransaction(wg sync.WaitGroup, transaction *types.Transaction, client *ethclient.Client) {
	defer wg.Done()
	txReceipt, err := eth.GetTransactionReceiptByHash(client, transaction.Hash())
	if err != nil {
		log.Errorf("GetTransactionReceiptByHash failed, err:%+v", err)
		return
	}
	ethValue := types2.ConvertWeiBitIntToEthDecimal(transaction.Value())
	data := transaction.Data()
	if len(data) == 0 && ethValue.IsPositive() {
		log.Infof("EventType:%s | txHash:%s | to:%s | amount:%v", "TRANSFER", transaction.Hash().Hex(), transaction.To().Hex(), ethValue)
		return
	}
	contractAddr := txReceipt.ContractAddress.Hex() // contract addr
	if transaction.To() == nil {
		// 可能是合约部署
		log.Infof("EventType:%s | ContractAddr:%s | transaction:%+v", "CONTRACT-DEPLOY", contractAddr, transaction)
		return
	}
	switch {
	case IsERC721Contract(common.HexToAddress(contractAddr), client) == true:
		log.Infof("EventType:%s | ContractAddr:%s | logs:%+v", "ERC721", contractAddr, DecodeTransferLog(txReceipt.Logs))
	case IsERC1155Contract(common.HexToAddress(contractAddr), client) == true:
		log.Infof("EventType:%s | ContractAddr:%s | logs:%+v | batchLogs", "ERC1155", contractAddr, DecodeTransferSingleLog(txReceipt.Logs), DecodeTransferBatchLog(txReceipt.Logs))
	}
}

func DecodeTransferLog(logs []*types.Log) []LogTransfer {
	var transferEvents []LogTransfer
	var transferEvent LogTransfer

	transferEventHash := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

	for _, vLog := range logs {
		if strings.Compare(vLog.Topics[0].Hex(), transferEventHash.Hex()) == 0 && len(vLog.Topics) >= 4 {
			func() {
				transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
				transferEvent.TokenId = vLog.Topics[3].Big()
				transferEvents = append(transferEvents, transferEvent)
			}()
		}
	}

	return transferEvents
}

func DecodeTransferSingleLog(logs []*types.Log) []LogTransferSingle {
	var transferEvents []LogTransferSingle
	var transferEvent LogTransferSingle

	transferEventHash := crypto.Keccak256Hash([]byte("TransferSingle(address,address,address,uint256,uint256)"))
	contractAbi, err := abi.JSON(strings.NewReader(string(contracts.ContractsABI)))

	for _, vLog := range logs {
		if strings.Compare(vLog.Topics[0].Hex(), transferEventHash.Hex()) == 0 && len(vLog.Topics) >= 4 {
			func() {
				err = contractAbi.UnpackIntoInterface(&transferEvent, "TransferSingle", vLog.Data)
				if err != nil {
					log.Fatal("fatal")
				}
				transferEvent.Operator = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.From = common.HexToAddress(vLog.Topics[2].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[3].Hex())

				transferEvents = append(transferEvents, transferEvent)
			}()
		}
	}

	return transferEvents
}

func DecodeTransferBatchLog(logs []*types.Log) []LogTransferBatch {
	var transferEvents []LogTransferBatch
	var transferEvent LogTransferBatch

	transferEventHash := crypto.Keccak256Hash([]byte("TransferBatch(address,address,address,uint256[],uint256[])"))
	contractAbi, err := abi.JSON(strings.NewReader(string(contracts.ContractsABI)))

	for _, vLog := range logs {
		if strings.Compare(vLog.Topics[0].Hex(), transferEventHash.Hex()) == 0 && len(vLog.Topics) >= 4 {
			func() {
				err = contractAbi.UnpackIntoInterface(&transferEvent, "TransferBatch", vLog.Data)
				if err != nil {
					log.Fatal("fatal")
				}
				transferEvent.Operator = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.From = common.HexToAddress(vLog.Topics[2].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[3].Hex())

				transferEvents = append(transferEvents, transferEvent)
			}()
		}
	}

	return transferEvents
}
