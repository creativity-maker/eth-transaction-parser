# eth-transaction-parser
eth-transaction-parser is an extensible Ethereum transaction parser.

## Feature
* transfer,approve
* Dex Protocol: Uniswap (TODO)


## Quickstart
Install eth-transaction-parser:
```shell
git clone https://github.com/ares0x/eth-transaction-parser.git
```
Configuration
Modify the PROVIDER in cmd/main.go to the URL with the API key you have obtained from [Infura](https://www.infura.io/). Of course, you can also utilize platforms like [quicknode](https://www.quicknode.com/), [alchemy](https://www.alchemy.com/)or set up your own node.

Running
```shell
cd eth-transaction-parser
go mod tidy
go run cmd/main.go
```