package main

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func newRPCClient(ctx context.Context, rpcURL string) *rpc.Client {
	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	return rpcClient
}

func newClient(rpcURL string) *ethclient.Client {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		panic(err)
	}
	return client
}
