package program

import (
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
)

var (
	Parsers = make(map[solana.PublicKey]Parser)
	Name2Id = make(map[string]solana.PublicKey)
	Id2Name = make(map[solana.PublicKey]string)
	Id2Type = make(map[solana.PublicKey]string)
)

const (
	Token      = "Token"
	Swap       = "Swap"
	StableSwap = "StableSwap"
	OrderBook  = "OrderBook"
)

type Parser func(in *types.Instruction, meta *types.Meta) error

func RegisterParser(program solana.PublicKey, name string, t string, p Parser) {
	Parsers[program] = p
	Name2Id[name] = program
	Id2Name[program] = name
	Id2Type[program] = t
}

func Parse(in *types.Instruction, meta *types.Meta) error {
	programId := in.Instruction.ProgramId
	parser, ok := Parsers[programId]
	if !ok {
		return nil
	}
	return parser(in, meta)
}
