package pump

import (
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/pumpfun"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(pumpfun.ProgramID, ProgramParser)
	RegisterParser(uint64(pumpfun.Instruction_Create.Uint32()), ParseInitializePool)
	RegisterParser(uint64(pumpfun.Instruction_Buy.Uint32()), ParseSwap)
	RegisterParser(uint64(pumpfun.Instruction_Sell.Uint32()), ParseSwapV2)
	RegisterParser(uint64(pumpfun.Instruction_Withdraw.Uint32()), ParseIncreaseLiquidity)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := pumpfun.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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
func ParseInitializePool(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Swap
func ParseSwap(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// SwapV2
func ParseSwapV2(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// IncreaseLiquidity
func ParseIncreaseLiquidity(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Default
func ParseDefault(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
