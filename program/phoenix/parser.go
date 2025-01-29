package phoenix_v1

import (
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/phoenix_v1"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(phoenix_v1.ProgramID, ProgramParser)
	RegisterParser(uint64(phoenix_v1.Instruction_InitializeMarket.Uint32()), ParseInitializePool)
	RegisterParser(uint64(phoenix_v1.Instruction_Swap.Uint32()), ParseSwap)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := phoenix_v1.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
	if err != nil {
		return
	}
	id := uint64(inst.TypeID.Uint32())
	parser, ok := Parsers[id]
	if !ok {
		return
	}
	parser(inst, in, meta)
}

// InitializePool
func ParseInitializePool(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Swap
func ParseSwap(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Default
func ParseDefault(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *phoenix_v1.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
