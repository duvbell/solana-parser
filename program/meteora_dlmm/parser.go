package meteora_dlmm

import (
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/meteora_dlmm"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(meteora_dlmm.ProgramID, ProgramParser)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializeLbPair.Uint32()), ParseInitializeLbPair)
	RegisterParser(uint64(meteora_dlmm.Instruction_AddLiquidity.Uint32()), ParseAddLiquidity)
	RegisterParser(uint64(meteora_dlmm.Instruction_AddLiquidityByStrategy.Uint32()), ParseAddLiquidityByStrategy)
	RegisterParser(uint64(meteora_dlmm.Instruction_RemoveLiquidity.Uint32()), ParseRemoveLiquidity)
	RegisterParser(uint64(meteora_dlmm.Instruction_RemoveLiquidityByRange.Uint32()), ParseRemoveLiquidityByRange)
	RegisterParser(uint64(meteora_dlmm.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(meteora_dlmm.Instruction_SwapExactOut.Uint32()), ParseSwapExactOut)
	RegisterParser(uint64(meteora_dlmm.Instruction_ClaimFee.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_ClaimReward.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_ClosePosition.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_ClosePresetParameter.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_FundReward.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_GoToABin.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_IncreaseOracleLength.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializeBinArray.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializeBinArrayBitmapExtension.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializePosition.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializePositionByOperator.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializePositionPda.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializeReward.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_MigrateBinArray.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_MigratePosition.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializePresetParameter.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_SetActivationPoint.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_SetPreActivationDuration.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_SetPreActivationSwapAddress.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_TogglePairStatus.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_UpdatePositionOperator.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_UpdateFeesAndRewards.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_UpdateFeeParameters.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_UpdateRewardDuration.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_UpdateRewardFunder.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_WithdrawProtocolFee.Uint32()), ParseDefault)
	RegisterParser(uint64(meteora_dlmm.Instruction_WithdrawIneligibleReward.Uint32()), ParseDefault)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := meteora_dlmm.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

// InitializeLbPair
func ParseInitializeLbPair(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// AddLiquidity
func ParseAddLiquidity(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_dlmm.AddLiquidity)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetLbPairAccount().PublicKey,
		User:           inst1.GetSenderAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	panic("not supported")
	in.Event = []interface{}{addLiquidity}
}

// AddLiquidityByStrategy
func ParseAddLiquidityByStrategy(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_dlmm.AddLiquidityByStrategy)
	transfers := in.FindChildrenWithTransfer()
	if len(transfers) != 2 {
		panic("not supported")
	}
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetLbPairAccount().PublicKey,
		User:           inst1.GetSenderAccount().PublicKey,
		TokenATransfer: transfers[0],
		TokenBTransfer: transfers[1],
	}
	in.Event = []interface{}{addLiquidity}
}

// RemoveLiquidity
func ParseRemoveLiquidity(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_dlmm.RemoveLiquidity)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	removeLiquidity := &types.RemoveLiquidity{
		Pool:           inst1.GetLbPairAccount().PublicKey,
		User:           inst1.GetSenderAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	panic("not supported")
	in.Event = []interface{}{removeLiquidity}
}

// RemoveLiquidityByRange
func ParseRemoveLiquidityByRange(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_dlmm.RemoveLiquidityByRange)
	transfers := in.FindChildrenWithTransfer()
	t1 := transfers[0]
	var t2 *types.Transfer
	if len(transfers) > 1 {
		t2 = transfers[1]
	}
	removeLiquidity := &types.RemoveLiquidity{
		Pool:           inst1.GetLbPairAccount().PublicKey,
		User:           inst1.GetSenderAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	in.Event = []interface{}{removeLiquidity}
}

// Swap
func ParseSwap(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_dlmm.Swap)
	// the first one is user deposit
	// the second is vault withdraw
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	swap := &types.Swap{
		Pool:           inst1.GetLbPairAccount().PublicKey,
		User:           inst1.GetUserAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	in.Event = []interface{}{swap}
}

// SwapExactOut
func ParseSwapExactOut(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_dlmm.SwapExactOut)
	// the first one is user deposit
	// the second is vault withdraw
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	swap := &types.Swap{
		Pool:           inst1.GetLbPairAccount().PublicKey,
		User:           inst1.GetUserAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	in.Event = []interface{}{swap}
}

// Default
func ParseDefault(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
