package raydium_clmm

import (
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/raydium_clmm"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(raydium_clmm.ProgramID, ProgramParser)
	RegisterParser(uint64(raydium_clmm.Instruction_CreatePool.Uint32()), ParseCreatePool)
	RegisterParser(uint64(raydium_clmm.Instruction_IncreaseLiquidityV2.Uint32()), ParseIncreaseLiquidityV2)
	RegisterParser(uint64(raydium_clmm.Instruction_DecreaseLiquidityV2.Uint32()), ParseDecreaseLiquidityV2)
	RegisterParser(uint64(raydium_clmm.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(raydium_clmm.Instruction_SwapV2.Uint32()), ParseSwapV2)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := raydium_clmm.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

// CreatePool
func ParseCreatePool(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_clmm.CreatePool)
	pool := &types.Pool{
		Hash:     inst1.GetPoolStateAccount().PublicKey,
		MintA:    inst1.GetTokenMint0Account().PublicKey,
		MintB:    inst1.GetTokenMint1Account().PublicKey,
		MintLp:   inst1.GetTokenVault1Account().PublicKey,
		VaultA:   inst1.GetTokenVault1Account().PublicKey,
		VaultB:   inst1.GetTokenVault1Account().PublicKey,
		ReserveA: 0,
		ReserveB: 0,
	}
	panic("not supported")
	in.Receipt = []interface{}{pool}
}

// IncreaseLiquidityV2
func ParseIncreaseLiquidityV2(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_clmm.IncreaseLiquidityV2)
	transfers := in.FindChildrenWithTransfer()
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetPoolStateAccount().PublicKey,
		User:           inst1.Get(0).PublicKey,
		TokenATransfer: transfers[0],
		TokenBTransfer: transfers[1],
	}
	in.Event = []interface{}{addLiquidity}
}

// DecreaseLiquidityV2
func ParseDecreaseLiquidityV2(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_clmm.DecreaseLiquidityV2)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	//
	removeLiquidity := &types.RemoveLiquidity{
		Pool:           inst1.GetPoolStateAccount().PublicKey,
		User:           inst1.Get(0).PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	panic("not supported")
	in.Event = []interface{}{removeLiquidity}
}

// Swap
func ParseSwap(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_clmm.Swap)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	//
	swap := &types.Swap{
		Pool:           inst1.GetPoolStateAccount().PublicKey,
		User:           inst1.GetPayerAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	in.Event = []interface{}{swap}
}

// SwapV2
func ParseSwapV2(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_clmm.SwapV2)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	swap := &types.Swap{
		Pool:           inst1.GetPoolStateAccount().PublicKey,
		User:           inst1.Get(0).PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	in.Event = []interface{}{swap}
}

// Default
func ParseDefault(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
