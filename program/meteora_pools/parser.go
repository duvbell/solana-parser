package meteora_pools

import (
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/meteora_pools"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(meteora_pools.ProgramID, ProgramParser)
	RegisterParser(uint64(meteora_pools.Instruction_InitializePermissionlessPool.Uint32()), ParseInitializePermissionlessPool)
	RegisterParser(uint64(meteora_pools.Instruction_AddBalanceLiquidity.Uint32()), ParseAddBalanceLiquidity)
	RegisterParser(uint64(meteora_pools.Instruction_RemoveBalanceLiquidity.Uint32()), ParseRemoveBalanceLiquidity)
	RegisterParser(uint64(meteora_pools.Instruction_Swap.Uint32()), ParseSwap)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := meteora_pools.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

// InitializePermissionlessPool
func ParseInitializePermissionlessPool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// AddBalanceLiquidity
func ParseAddBalanceLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_pools.AddBalanceLiquidity)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetPoolAccount().PublicKey,
		User:           inst1.GetUserAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	panic("not supported")
	in.Event = []interface{}{addLiquidity}
}

// RemoveBalanceLiquidity
func ParseRemoveBalanceLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_pools.RemoveBalanceLiquidity)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	removeLiquidity := &types.RemoveLiquidity{
		Pool:           inst1.GetPoolAccount().PublicKey,
		User:           inst1.GetUserAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	panic("not supported")
	in.Event = []interface{}{removeLiquidity}
}

// Swap
func ParseSwap(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_pools.Swap)
	// the transfer is execute by vault deposit & withdraw & this first transfer is fee
	deposit := in.Children[0].Event[0].(*types.Transfer)
	withdraw := in.Children[1].Event[0].(*types.Transfer)
	swap := &types.Swap{
		Pool:           inst1.GetPoolAccount().PublicKey,
		User:           inst1.GetUserAccount().PublicKey,
		TokenATransfer: deposit,
		TokenBTransfer: withdraw,
	}
	in.Event = []interface{}{swap}
}

// Default
func ParseDefault(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
