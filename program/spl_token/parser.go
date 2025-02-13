package spl_token

import (
	"encoding/json"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
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
	program.RegisterParser(solana.TokenProgramID, "token", program.Token, ProgramParser)
	RegisterParser(types.CreateId([]byte("transfer")), ParseTransfer)
	RegisterParser(types.CreateId([]byte("transferChecked")), ParseTransfer)
	RegisterParser(types.CreateId([]byte("mintTo")), ParseMint)
	RegisterParser(types.CreateId([]byte("burn")), ParseBurn)
	RegisterParser(types.CreateId([]byte("initializeAccount")), ParseInitialize)
	RegisterParser(types.CreateId([]byte("initializeAccount3")), ParseInitialize)
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
	id := types.CreateId([]byte(instruction.T))
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
	Amount      uint64           `json:"amount,string"`
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
	amount := instruction.Amount
	if amount == 0 {
		amount = instruction.TokenAmount.Amount
	}
	tokenAccount := meta.TokenAccounts[from]
	transfer := &types.Transfer{
		Mint:   tokenAccount.Mint,
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
	meta.TokenAccounts[init.Account] = &types.TokenAccount{
		Owner:     &instruction.Owner,
		ProgramId: &in.Instruction.ProgramId,
		Mint:      instruction.Mint,
	}
	in.Event = []interface{}{init}
}
