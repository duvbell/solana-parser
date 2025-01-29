package whirlpool

import (
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/whirlpool"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(whirlpool.ProgramID, ProgramParser)
	RegisterParser(uint64(whirlpool.Instruction_InitializePool.Uint32()), ParseInitializePool)
	RegisterParser(uint64(whirlpool.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(whirlpool.Instruction_SwapV2.Uint32()), ParseSwapV2)
	RegisterParser(uint64(whirlpool.Instruction_IncreaseLiquidity.Uint32()), ParseIncreaseLiquidity)
	RegisterParser(uint64(whirlpool.Instruction_DecreaseLiquidity.Uint32()), ParseDecreaseLiquidity)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := whirlpool.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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
func ParseInitializePool(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Swap
func ParseSwap(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*whirlpool.Swap)
	// child 1 : transfer
	// child 2 : transfer
	transfers := in.FindChildrenWithTransfer()
	swap := &types.Swap{
		Pool:           inst1.GetWhirlpoolAccount().PublicKey,
		TokenATransfer: transfers[0],
		TokenBTransfer: transfers[1],
		User:           inst1.GetTokenOwnerAccountAAccount().PublicKey,
	}
	in.Event = []interface{}{swap}
}

// SwapV2
func ParseSwapV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*whirlpool.SwapV2)
	// child 1 : transfer
	// child 2 : transfer
	transfers := in.FindChildrenWithTransfer()
	swap := &types.Swap{
		Pool:           inst1.GetWhirlpoolAccount().PublicKey,
		TokenATransfer: transfers[0],
		TokenBTransfer: transfers[1],
		User:           inst1.GetTokenOwnerAccountAAccount().PublicKey,
	}
	in.Event = []interface{}{swap}
}

// IncreaseLiquidity
func ParseIncreaseLiquidity(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	// child 1 : transfer
	// child 2 : transfer
	transfer1 := in.Children[0].Event[0]
	transfer2 := in.Children[0].Event[0]
	in.Event = []interface{}{transfer1, transfer2}
}

// DecreaseLiquidity
func ParseDecreaseLiquidity(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	// child 1 : transfer
	// child 2 : transfer
	transfer1 := in.Children[0].Event[0]
	transfer2 := in.Children[0].Event[0]
	in.Event = []interface{}{transfer1, transfer2}
}

// Default
func ParseDefault(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
