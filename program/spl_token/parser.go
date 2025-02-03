package spl_token

import (
	"encoding/json"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
	"math/big"
)

var (
	programId = solana.TokenProgramID
	Parsers   = make(map[uint64]Parser, 0)
)

type Parser func(in *types.Instruction, raw []byte, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(programId, ProgramParser)
	RegisterParser(new(big.Int).SetBytes([]byte("transfer")).Uint64(), ParseTransfer)
	RegisterParser(new(big.Int).SetBytes([]byte("transferChecked")).Uint64(), ParseTransfer)
	RegisterParser(new(big.Int).SetBytes([]byte("mintTo")).Uint64(), ParseMint)
	RegisterParser(new(big.Int).SetBytes([]byte("burn")).Uint64(), ParseBurn)
	RegisterParser(new(big.Int).SetBytes([]byte("initializeAccount")).Uint64(), ParseInitialize)
	RegisterParser(new(big.Int).SetBytes([]byte("initializeAccount3")).Uint64(), ParseInitialize)
}

type Instruction struct {
	T   string          `json:"type"`
	Raw json.RawMessage `json:"info"`
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	inJson, _ := in.Instruction.Parsed.MarshalJSON()
	var instruction Instruction
	err := json.Unmarshal(inJson, &instruction)
	if err != nil {
		return err
	}
	id := new(big.Int).SetBytes([]byte(instruction.T)).Uint64()
	parser, ok := Parsers[id]
	if !ok {
		return nil
	}
	parser(in, instruction.Raw, meta)
	return nil
}

// transfer instruction
type Transfer struct {
	Destination solana.PublicKey `json:"destination"`
	Lamports    uint64           `json:"amount,string"`
	Source      solana.PublicKey `json:"source"`
	Authority   solana.PublicKey `json:"authority"`
	Mint        solana.PublicKey `json:"mint"`
	TokenAmount struct {
		Amount   uint64 `json:"amount,string"`
		Decimals uint64 `json:"decimals"`
	} `json:"tokenAmount"`
}

func ParseTransfer(in *types.Instruction, raw []byte, meta *types.Meta) {
	instruction := &Transfer{}
	if err := json.Unmarshal(raw, instruction); err != nil {
		return
	}
	mint := meta.TokenMint[instruction.Source]
	from := instruction.Source
	/*
		if k, ok := meta.TokenOwner[instruction.Source]; ok {
			from = k
		}
	*/
	to := instruction.Destination
	/*
		if k, ok := meta.TokenOwner[instruction.Destination]; ok {
			to = k
		}
	*/
	amount := instruction.Lamports
	if amount == 0 {
		amount = instruction.TokenAmount.Amount
	}
	transfer := &types.Transfer{
		Mint:   mint,
		Amount: amount,
		From:   from,
		To:     to,
	}
	in.Event = []interface{}{transfer}
}

// mint instruction
type Mint struct {
	Account   solana.PublicKey `json:"account"`
	Amount    uint64           `json:"amount,string"`
	Authority solana.PublicKey `json:"mintAuthority"`
	Mint      solana.PublicKey `json:"mint"`
}

func ParseMint(in *types.Instruction, raw []byte, meta *types.Meta) {
	instruction := &Mint{}
	if err := json.Unmarshal(raw, instruction); err != nil {
		return
	}
	account := instruction.Account
	/*
		if k, ok := meta.TokenOwner[instruction.Account]; ok {
			account = k
		}
	*/
	mintTo := &types.MintTo{
		Mint:    instruction.Mint,
		Amount:  instruction.Amount,
		Account: account,
	}
	in.Event = []interface{}{mintTo}
}

// burn instruction
type Burn struct {
	Account   solana.PublicKey `json:"account"`
	Amount    uint64           `json:"amount,string"`
	Authority solana.PublicKey `json:"authority"`
	Mint      solana.PublicKey `json:"mint"`
}

func ParseBurn(in *types.Instruction, raw []byte, meta *types.Meta) {
	instruction := &Burn{}
	if err := json.Unmarshal(raw, instruction); err != nil {
		return
	}
	account := instruction.Account
	/*
		if k, ok := meta.TokenOwner[instruction.Account]; ok {
			account = k
		}
	*/
	burn := &types.Burn{
		Mint:    instruction.Mint,
		Amount:  instruction.Amount,
		Account: account,
	}
	in.Event = []interface{}{burn}
}

// Initialize instruction
type Initialize struct {
	Account solana.PublicKey `json:"account"`
	Owner   solana.PublicKey `json:"owner"`
	Mint    solana.PublicKey `json:"mint"`
}

func ParseInitialize(in *types.Instruction, raw []byte, meta *types.Meta) {
	instruction := &Initialize{}
	if err := json.Unmarshal(raw, instruction); err != nil {
		return
	}
	init := &types.Initialize{
		Mint:    instruction.Mint,
		Account: instruction.Account,
		Owner:   instruction.Owner,
	}
	// update token owner & mint by spl token instructions
	meta.TokenOwner[init.Account] = init.Owner
	meta.TokenMint[init.Account] = init.Mint
	in.Event = []interface{}{init}
}
