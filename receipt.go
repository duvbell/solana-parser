package solanaparser

type Mint struct {
	Hash        string
	Owner       string
	Name        string
	Symbol      string
	Decimal     uint64
	TotalSupply uint64
}

type Token struct {
	User  string
	Mint  string
	Owner string
}

type Dex struct {
	Id   string
	Name string
}

type Pool struct {
	Hash     string
	MintA    string
	MintB    string
	MintLp   string
	VaultA   string
	VaultB   string
	VaultLp  string
	ReserveA uint64
	ReserveB uint64
}
