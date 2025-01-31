package meteora_dlmm

import (
	"github.com/blockchain-develop/solana-parser/log"
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
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializePermissionLbPair.Uint32()), ParseInitializePermissionLbPair)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializeCustomizablePermissionlessLbPair.Uint32()), ParseInitializeCustomizablePermissionlessLbPair)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializeBinArrayBitmapExtension.Uint32()), ParseInitializeBinArrayBitmapExtension)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializeBinArray.Uint32()), ParseInitializeBinArray)
	RegisterParser(uint64(meteora_dlmm.Instruction_AddLiquidity.Uint32()), ParseAddLiquidity)
	RegisterParser(uint64(meteora_dlmm.Instruction_AddLiquidityByWeight.Uint32()), ParseAddLiquidityByWeight)
	RegisterParser(uint64(meteora_dlmm.Instruction_AddLiquidityByStrategy.Uint32()), ParseAddLiquidityByStrategy)
	RegisterParser(uint64(meteora_dlmm.Instruction_AddLiquidityByStrategyOneSide.Uint32()), ParseAddLiquidityByStrategyOneSide)
	RegisterParser(uint64(meteora_dlmm.Instruction_AddLiquidityOneSide.Uint32()), ParseAddLiquidityOneSide)
	RegisterParser(uint64(meteora_dlmm.Instruction_RemoveLiquidity.Uint32()), ParseRemoveLiquidity)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializePosition.Uint32()), ParseInitializePosition)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializePositionPda.Uint32()), ParseInitializePositionPda)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializePositionByOperator.Uint32()), ParseInitializePositionByOperator)
	RegisterParser(uint64(meteora_dlmm.Instruction_UpdatePositionOperator.Uint32()), ParseUpdatePositionOperator)
	RegisterParser(uint64(meteora_dlmm.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(meteora_dlmm.Instruction_SwapExactOut.Uint32()), ParseSwapExactOut)
	RegisterParser(uint64(meteora_dlmm.Instruction_SwapWithPriceImpact.Uint32()), ParseSwapWithPriceImpact)
	RegisterParser(uint64(meteora_dlmm.Instruction_WithdrawProtocolFee.Uint32()), ParseWithdrawProtocolFee)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializeReward.Uint32()), ParseInitializeReward)
	RegisterParser(uint64(meteora_dlmm.Instruction_FundReward.Uint32()), ParseFundReward)
	RegisterParser(uint64(meteora_dlmm.Instruction_UpdateRewardFunder.Uint32()), ParseUpdateRewardFunder)
	RegisterParser(uint64(meteora_dlmm.Instruction_UpdateRewardDuration.Uint32()), ParseUpdateRewardDuration)
	RegisterParser(uint64(meteora_dlmm.Instruction_ClaimReward.Uint32()), ParseClaimReward)
	RegisterParser(uint64(meteora_dlmm.Instruction_ClaimFee.Uint32()), ParseClaimFee)
	RegisterParser(uint64(meteora_dlmm.Instruction_ClosePosition.Uint32()), ParseClosePosition)
	RegisterParser(uint64(meteora_dlmm.Instruction_UpdateFeeParameters.Uint32()), ParseUpdateFeeParameters)
	RegisterParser(uint64(meteora_dlmm.Instruction_IncreaseOracleLength.Uint32()), ParseIncreaseOracleLength)
	RegisterParser(uint64(meteora_dlmm.Instruction_InitializePresetParameter.Uint32()), ParseInitializePresetParameter)
	RegisterParser(uint64(meteora_dlmm.Instruction_ClosePresetParameter.Uint32()), ParseClosePresetParameter)
	RegisterParser(uint64(meteora_dlmm.Instruction_RemoveAllLiquidity.Uint32()), ParseRemoveAllLiquidity)
	RegisterParser(uint64(meteora_dlmm.Instruction_TogglePairStatus.Uint32()), ParseTogglePairStatus)
	RegisterParser(uint64(meteora_dlmm.Instruction_MigratePosition.Uint32()), ParseMigratePosition)
	RegisterParser(uint64(meteora_dlmm.Instruction_MigrateBinArray.Uint32()), ParseMigrateBinArray)
	RegisterParser(uint64(meteora_dlmm.Instruction_UpdateFeesAndRewards.Uint32()), ParseUpdateFeesAndRewards)
	RegisterParser(uint64(meteora_dlmm.Instruction_WithdrawIneligibleReward.Uint32()), ParseWithdrawIneligibleReward)
	RegisterParser(uint64(meteora_dlmm.Instruction_SetActivationPoint.Uint32()), ParseSetActivationPoint)
	RegisterParser(uint64(meteora_dlmm.Instruction_RemoveLiquidityByRange.Uint32()), ParseRemoveLiquidityByRange)
	RegisterParser(uint64(meteora_dlmm.Instruction_AddLiquidityOneSidePrecise.Uint32()), ParseAddLiquidityOneSidePrecise)
	RegisterParser(uint64(meteora_dlmm.Instruction_GoToABin.Uint32()), ParseGoToABin)
	RegisterParser(uint64(meteora_dlmm.Instruction_SetPreActivationDuration.Uint32()), ParseSetPreActivationDuration)
	RegisterParser(uint64(meteora_dlmm.Instruction_SetPreActivationSwapAddress.Uint32()), ParseSetPreActivationSwapAddress)
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

func ParseInitializeLbPair(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	// create all accounts
}
func ParseInitializePermissionLbPair(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize permission lb pair", "program", meteora_dlmm.ProgramName)
}
func ParseInitializeCustomizablePermissionlessLbPair(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize customizable permissionless lb pair", "program", meteora_dlmm.ProgramName)
}
func ParseInitializeBinArrayBitmapExtension(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitializeBinArray(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
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
	in.Event = []interface{}{addLiquidity}
}
func ParseAddLiquidityByWeight(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_dlmm.AddLiquidityByWeight)
	transfers := in.FindChildrenTransfers()
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
func ParseAddLiquidityByStrategy(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_dlmm.AddLiquidityByStrategy)
	transfers := in.FindChildrenTransfers()
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
func ParseAddLiquidityByStrategyOneSide(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse add liquidity by strategy one-side", "program", meteora_dlmm.ProgramName)
}
func ParseAddLiquidityOneSide(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse add liquidity one-side", "program", meteora_dlmm.ProgramName)
}
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
	in.Event = []interface{}{removeLiquidity}
}
func ParseInitializePosition(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	// only create accounts
}
func ParseInitializePositionPda(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitializePositionByOperator(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdatePositionOperator(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
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
func ParseSwapWithPriceImpact(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse swap with price impact", "program", meteora_dlmm.ProgramName)
}
func ParseWithdrawProtocolFee(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitializeReward(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseFundReward(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdateRewardFunder(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdateRewardDuration(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseClaimReward(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseClaimFee(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseClosePosition(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdateFeeParameters(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseIncreaseOracleLength(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitializePresetParameter(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseClosePresetParameter(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseRemoveAllLiquidity(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse remove all liquidity", "program", meteora_dlmm.ProgramName)
}
func ParseTogglePairStatus(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseMigratePosition(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseMigrateBinArray(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdateFeesAndRewards(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseWithdrawIneligibleReward(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetActivationPoint(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseRemoveLiquidityByRange(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_dlmm.RemoveLiquidityByRange)
	transfers := in.FindChildrenTransfers()
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
func ParseAddLiquidityOneSidePrecise(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse add liquidity one-side precise", "program", meteora_dlmm.ProgramName)
}
func ParseGoToABin(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetPreActivationDuration(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetPreActivationSwapAddress(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
}

// Default
func ParseDefault(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *meteora_dlmm.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
