package program

import (
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
)

var (
	Parsers = make(map[solana.PublicKey]Parser, 0)
	Name2Id = make(map[string]solana.PublicKey, 0)
	Id2Name = make(map[solana.PublicKey]string, 0)
)

type Parser func(in *types.Instruction, meta *types.Meta) error

func RegisterParser(program solana.PublicKey, name string, p Parser) {
	Parsers[program] = p
	Name2Id[name] = program
	Id2Name[program] = name
}

func Parse(in *types.Instruction, meta *types.Meta) error {
	programId := in.Instruction.ProgramId
	parser, ok := Parsers[programId]
	if !ok {
		return nil
	}
	return parser(in, meta)
}
