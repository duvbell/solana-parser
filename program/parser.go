package program

import (
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
)

var (
	Parsers     = make(map[solana.PublicKey]Parser)
	Name2Id     = make(map[string]solana.PublicKey)
	Id2Name     = make(map[solana.PublicKey]string)
	Id2Type     = make(map[solana.PublicKey]string)
	Id2Priority = make(map[solana.PublicKey]int)
)

const (
	Token      = "Token"
	Swap       = "Swap"
	StableSwap = "StableSwap"
	OrderBook  = "OrderBook"
)

type Parser func(in *types.Instruction, meta *types.Meta) error

func RegisterParser(program solana.PublicKey, name string, t string, priority int, p Parser) {
	Parsers[program] = p
	Id2Priority[program] = priority
	Name2Id[name] = program
	Id2Name[program] = name
	Id2Type[program] = t
}

func RemoveParser(program solana.PublicKey) {
	Parsers[program] = nil
}

func FilterInstruction(in *solana.CompiledInstruction, meta *types.Meta) *types.Instruction {
	programId := meta.Accounts[in.ProgramIDIndex].PublicKey
	accountMetas := make([]*solana.AccountMeta, 0)
	for _, accountIndex := range in.Accounts {
		accountMetas = append(accountMetas, meta.Accounts[accountIndex])
	}
	return &types.Instruction{
		RawInstruction: &solana.GenericInstruction{
			AccountValues: accountMetas,
			ProgID:        programId,
			DataBytes:     in.Data,
		},
	}
}

func Parse(in *types.Instruction, meta *types.Meta) error {
	programId := in.RawInstruction.ProgID
	parser, ok := Parsers[programId]
	if !ok {
		return nil
	}
	return parser(in, meta)
}
