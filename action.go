package solanaparser

import "github.com/shopspring/decimal"

type Trade struct {
	Pool         string
	User         string
	Type         string
	TokenAAmount decimal.Decimal
	TokenBAmount decimal.Decimal
}

type CreatePool struct {
	Pool      string
	TokenA    string
	TokenB    string
	TokenLP   string
	AccountA  string
	AccountB  string
	AccountLP string
	User      string
}

type AddLiquidity struct {
	Pool           string
	TokenATransfer *Transfer
	TokenBTransfer *Transfer
	TokenLpMint    *MintTo
	User           string
}

type RemoveLiquidity struct {
	Pool           string
	TokenATransfer *Transfer
	TokenBTransfer *Transfer
	TokenLpBurn    *Burn
	User           string
}

type Swap struct {
	Pool           string
	TokenATransfer *Transfer
	TokenBTransfer *Transfer
	User           string
}

type Transfer struct {
	Mint   string
	From   string
	To     string
	Amount uint64
}

type MintTo struct {
	Mint    string
	Account string
	Amount  uint64
}

type Burn struct {
	Mint    string
	Account string
	Amount  uint64
}
