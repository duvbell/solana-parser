package transaction

import "github.com/shopspring/decimal"

const (
	CreatePool      = "create_pool"
	AddLiquidity    = "add_liquidity"
	RemoveLiquidity = "remove_liquidity"
	Swap            = "swap"
)

type Trade struct {
	Pool         string
	User         string
	Type         string
	TokenAAmount decimal.Decimal
	TokenBAmount decimal.Decimal
}

type Transfer struct {
	Mint   string
	From   string
	To     string
	Amount uint64
}
