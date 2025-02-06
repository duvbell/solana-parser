package types

import "github.com/gagliardetto/solana-go"

type CreatePool struct {
	Dex     solana.PublicKey
	Pool    solana.PublicKey
	User    solana.PublicKey
	TokenA  solana.PublicKey
	TokenB  solana.PublicKey
	TokenLP solana.PublicKey
	VaultA  solana.PublicKey
	VaultB  solana.PublicKey
	VaultLP solana.PublicKey
}

type AddLiquidity struct {
	Dex            solana.PublicKey
	Pool           solana.PublicKey
	User           solana.PublicKey
	TokenATransfer *Transfer
	TokenBTransfer *Transfer
	TokenLpMint    *MintTo
}

type RemoveLiquidity struct {
	Dex            solana.PublicKey
	Pool           solana.PublicKey
	User           solana.PublicKey
	TokenATransfer *Transfer
	TokenBTransfer *Transfer
	TokenLpBurn    *Burn
}

type Swap struct {
	Dex            solana.PublicKey
	Pool           solana.PublicKey
	User           solana.PublicKey
	InputTransfer  *Transfer
	OutputTransfer *Transfer
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
