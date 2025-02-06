package system

import (
	"encoding/json"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
	"math/big"
)

var (
	programId = solana.SystemProgramID
	Parsers   = make(map[uint64]Parser, 0)
)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

type Parser func(in *types.Instruction, raw []byte, meta *types.Meta)

func init() {
	program.RegisterParser(programId, "system", ProgramParser)
	RegisterParser(new(big.Int).SetBytes([]byte("transfer")).Uint64(), ParseTransfer)
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

type Transfer struct {
	Destination solana.PublicKey `json:"destination"`
	Lamports    uint64           `json:"lamports"`
	Source      solana.PublicKey `json:"source"`
}

func ParseTransfer(in *types.Instruction, raw []byte, meta *types.Meta) {
	instruction := &Transfer{}
	if err := json.Unmarshal(raw, instruction); err != nil {
		return
	}
	transfer := &types.Transfer{
		Mint:   solana.MustPublicKeyFromBase58("11111111111111111111111111111111"),
		Amount: instruction.Lamports,
		From:   instruction.Source,
		To:     instruction.Destination,
	}
	in.Event = []interface{}{transfer}
}
