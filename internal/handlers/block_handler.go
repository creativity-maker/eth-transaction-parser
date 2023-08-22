package handlers

import (
	"context"
	"eth-transaction-parser/pkg/eth"
	"eth-transaction-parser/pkg/log"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"sync"
	"time"
)

type BlockHandler struct {
	Name        string
	BlockHeight uint64
	Client      *ethclient.Client
	Workers     int
}

func NewBlockHandler(client *ethclient.Client) *BlockHandler {
	return &BlockHandler{
		Name:        "BLOCK_HANDLER",
		BlockHeight: 0,
		Client:      client,
		Workers:     300,
	}
}

func (h *BlockHandler) Run() {
	ticker := time.NewTicker(time.Second * 10)
	for _ = range ticker.C {
		h.RotatingThroughBlocks()
	}
}

func (h *BlockHandler) RotatingThroughBlocks() {
	height, err := eth.GetBlockHeight(h.Client)
	if err != nil {
		return
	}
	if h.BlockHeight != 0 && h.BlockHeight >= height {
		return
	} else {
		h.BlockHeight = height
	}
	block, err := h.Client.BlockByNumber(context.TODO(), big.NewInt(int64(h.BlockHeight)))
	if err != nil {
		log.Fatal("fatal")
	}
	//log.Infof("block:%+v", block)
	transactions := block.Transactions()
	//ch := make(chan interface{}, h.Workers)
	var wg sync.WaitGroup
	for _, transaction := range transactions {
		wg.Add(1)
		go ParseTransaction(wg, transaction, h.Client)
	}
	wg.Wait()
	h.BlockHeight += 1
}
