package stable_swap

import (
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/stable_swap"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(stable_swap.ProgramID, ProgramParser)
	RegisterParser(uint64(stable_swap.Instruction_Initialize.Uint32()), ParseInitialize)
	RegisterParser(uint64(stable_swap.Instruction_Deposit.Uint32()), ParseDeposit)
	RegisterParser(uint64(stable_swap.Instruction_Withdraw.Uint32()), ParseWithdraw)
	RegisterParser(uint64(stable_swap.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(stable_swap.Instruction_SwapV2.Uint32()), ParseSwapV2)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := stable_swap.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

// Initialize
func ParseInitialize(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Deposit
func ParseDeposit(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*stable_swap.Deposit)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	user := inst1.GetUserPoolTokenAccount().PublicKey
	if owner, ok := meta.TokenOwner[user]; ok {
		user = owner
	}
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetPoolAccount().PublicKey,
		User:           user,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	panic("not supported")
	in.Event = []interface{}{addLiquidity}
}

// Withdraw
func ParseWithdraw(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*stable_swap.Withdraw)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	user := inst1.GetUserPoolTokenAccount().PublicKey
	if owner, ok := meta.TokenOwner[user]; ok {
		user = owner
	}
	removeLiquidity := &types.RemoveLiquidity{
		Pool:           inst1.GetPoolAccount().PublicKey,
		User:           user,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	panic("not supported")
	in.Event = []interface{}{removeLiquidity}
}

// Swap
func ParseSwap(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*stable_swap.Swap)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	user := inst1.GetUserTokenInAccount().PublicKey
	if owner, ok := meta.TokenOwner[user]; ok {
		user = owner
	}
	swap := &types.Swap{
		Pool:           inst1.GetPoolAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
		User:           user,
	}
	in.Event = []interface{}{swap}
}

// SwapV2
func ParseSwapV2(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*stable_swap.SwapV2)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	user := inst1.GetUserTokenInAccount().PublicKey
	if owner, ok := meta.TokenOwner[user]; ok {
		user = owner
	}
	swap := &types.Swap{
		Pool:           inst1.GetPoolAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
		User:           user,
	}
	in.Event = []interface{}{swap}
}

// Default
func ParseDefault(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
