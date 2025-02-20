package types

import "github.com/gagliardetto/solana-go"

type Mint struct {
	Hash        string
	Owner       string
	Name        string
	Symbol      string
	Decimal     uint64
	TotalSupply uint64
}

type Token struct {
	User  solana.PublicKey
	Mint  solana.PublicKey
	Owner solana.PublicKey
}

type Dex struct {
	Id   solana.PublicKey
	Name string
}

type Pool struct {
	Hash     solana.PublicKey
	MintA    solana.PublicKey
	MintB    solana.PublicKey
	MintLp   solana.PublicKey
	VaultA   solana.PublicKey
	VaultB   solana.PublicKey
	ReserveA uint64
	ReserveB uint64
}

type MemeCreateEvent struct {
	Name         string
	Symbol       string
	Uri          string
	Mint         solana.PublicKey
	BondingCurve solana.PublicKey
	User         solana.PublicKey
}

type MemeBuyEvent struct {
	Mint                 solana.PublicKey
	SolAmount            uint64
	TokenAmount          uint64
	IsBuy                bool
	User                 solana.PublicKey
	Timestamp            int64
	VirtualSolReserves   uint64
	VirtualTokenReserves uint64
}

type MemeSellEvent struct {
	Mint                 solana.PublicKey
	SolAmount            uint64
	TokenAmount          uint64
	IsBuy                bool
	User                 solana.PublicKey
	Timestamp            int64
	VirtualSolReserves   uint64
	VirtualTokenReserves uint64
}
