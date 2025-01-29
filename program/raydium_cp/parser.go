package raydium_cp

import (
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/raydium_cp"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(raydium_cp.ProgramID, ProgramParser)
	RegisterParser(uint64(raydium_cp.Instruction_Initialize.Uint32()), ParseInitializePool)
	RegisterParser(uint64(raydium_cp.Instruction_Deposit.Uint32()), ParseSwap)
	RegisterParser(uint64(raydium_cp.Instruction_Withdraw.Uint32()), ParseSwapV2)
	RegisterParser(uint64(raydium_cp.Instruction_SwapBaseInput.Uint32()), ParseIncreaseLiquidity)
	RegisterParser(uint64(raydium_cp.Instruction_SwapBaseOutput.Uint32()), ParseDecreaseLiquidity)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := raydium_cp.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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
func ParseInitializePool(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Swap
func ParseSwap(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// SwapV2
func ParseSwapV2(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// IncreaseLiquidity
func ParseIncreaseLiquidity(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// DecreaseLiquidity
func ParseDecreaseLiquidity(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Default
func ParseDefault(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
