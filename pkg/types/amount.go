package types

import (
	"github.com/shopspring/decimal"
	"math/big"
)

const WEI = 1000000000000000000

func ConvertWeiDecimalToEth(amount decimal.Decimal) decimal.Decimal {
	return amount.Div(decimal.NewFromFloat(WEI))
}

func ConvertWeiBitIntToEthDecimal(amount *big.Int) decimal.Decimal {
	return decimal.NewFromBigInt(amount, 0).Div(decimal.NewFromFloat(WEI))
}
