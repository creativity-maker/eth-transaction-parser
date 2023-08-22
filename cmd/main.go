package main

import (
	"eth-transaction-parser/internal/handlers"
	"eth-transaction-parser/pkg/log"
)

const (
	PROVIDER = ""
)

func main() {
	log.Init(nil)
	client := newClient(PROVIDER)
	handlers.NewBlockHandler(client).Run()
}
