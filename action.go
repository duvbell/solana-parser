package solanaparser

import "github.com/gagliardetto/solana-go"

type CreatePool struct {
	Pool      solana.PublicKey
	TokenA    solana.PublicKey
	TokenB    solana.PublicKey
	TokenLP   solana.PublicKey
	AccountA  solana.PublicKey
	AccountB  solana.PublicKey
	AccountLP solana.PublicKey
	User      solana.PublicKey
}

type AddLiquidity struct {
	Pool           solana.PublicKey
	TokenATransfer *Transfer
	TokenBTransfer *Transfer
	TokenLpMint    *MintTo
	User           solana.PublicKey
}

type RemoveLiquidity struct {
	Pool           solana.PublicKey
	TokenATransfer *Transfer
	TokenBTransfer *Transfer
	TokenLpBurn    *Burn
	User           solana.PublicKey
}

type Swap struct {
	Pool           solana.PublicKey
	TokenATransfer *Transfer
	TokenBTransfer *Transfer
	User           solana.PublicKey
}

type Transfer struct {
	Mint   solana.PublicKey
	From   solana.PublicKey
	To     solana.PublicKey
	Amount uint64
}

type MintTo struct {
	Mint    solana.PublicKey
	Account solana.PublicKey
	Amount  uint64
}

type Burn struct {
	Mint    solana.PublicKey
	Account solana.PublicKey
	Amount  uint64
}

type Initialize struct {
	Account solana.PublicKey
	Owner   solana.PublicKey
	Mint    solana.PublicKey
}
