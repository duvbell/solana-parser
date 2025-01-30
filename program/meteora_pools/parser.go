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
	RegisterParser(uint64(meteora_pools.Instruction_InitializePermissionedPool.Uint32()), ParseInitializePermissionedPool)
	RegisterParser(uint64(meteora_pools.Instruction_InitializePermissionlessPool.Uint32()), ParseInitializePermissionlessPool)
	RegisterParser(uint64(meteora_pools.Instruction_InitializePermissionlessPoolWithFeeTier.Uint32()), ParseInitializePermissionlessPoolWithFeeTier)
	RegisterParser(uint64(meteora_pools.Instruction_EnableOrDisablePool.Uint32()), ParseEnableOrDisablePool)
	RegisterParser(uint64(meteora_pools.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(meteora_pools.Instruction_RemoveLiquiditySingleSide.Uint32()), ParseRemoveLiquiditySingleSide)
	RegisterParser(uint64(meteora_pools.Instruction_AddImbalanceLiquidity.Uint32()), ParseAddImbalanceLiquidity)
	RegisterParser(uint64(meteora_pools.Instruction_RemoveBalanceLiquidity.Uint32()), ParseRemoveBalanceLiquidity)
	RegisterParser(uint64(meteora_pools.Instruction_AddBalanceLiquidity.Uint32()), ParseAddBalanceLiquidity)
	RegisterParser(uint64(meteora_pools.Instruction_SetPoolFees.Uint32()), ParseSetPoolFees)
	RegisterParser(uint64(meteora_pools.Instruction_OverrideCurveParam.Uint32()), ParseOverrideCurveParam)
	RegisterParser(uint64(meteora_pools.Instruction_GetPoolInfo.Uint32()), ParseGetPoolInfo)
	RegisterParser(uint64(meteora_pools.Instruction_BootstrapLiquidity.Uint32()), ParseBootstrapLiquidity)
	RegisterParser(uint64(meteora_pools.Instruction_CreateMintMetadata.Uint32()), ParseCreateMintMetadata)
	RegisterParser(uint64(meteora_pools.Instruction_CreateLockEscrow.Uint32()), ParseCreateLockEscrow)
	RegisterParser(uint64(meteora_pools.Instruction_Lock.Uint32()), ParseLock)
	RegisterParser(uint64(meteora_pools.Instruction_ClaimFee.Uint32()), ParseClaimFee)
	RegisterParser(uint64(meteora_pools.Instruction_CreateConfig.Uint32()), ParseCreateConfig)
	RegisterParser(uint64(meteora_pools.Instruction_CloseConfig.Uint32()), ParseCloseConfig)
	RegisterParser(uint64(meteora_pools.Instruction_InitializePermissionlessConstantProductPoolWithConfig.Uint32()), ParseInitializePermissionlessConstantProductPoolWithConfig)
	RegisterParser(uint64(meteora_pools.Instruction_InitializePermissionlessConstantProductPoolWithConfig2.Uint32()), ParseInitializePermissionlessConstantProductPoolWithConfig2)
	RegisterParser(uint64(meteora_pools.Instruction_InitializeCustomizablePermissionlessConstantProductPool.Uint32()), ParseInitializeCustomizablePermissionlessConstantProductPool)
	RegisterParser(uint64(meteora_pools.Instruction_UpdateActivationPoint.Uint32()), ParseUpdateActivationPoint)
	RegisterParser(uint64(meteora_pools.Instruction_WithdrawProtocolFees.Uint32()), ParseWithdrawProtocolFees)
	RegisterParser(uint64(meteora_pools.Instruction_SetWhitelistedVault.Uint32()), ParseSetWhitelistedVault)
	RegisterParser(uint64(meteora_pools.Instruction_PartnerClaimFee.Uint32()), ParsePartnerClaimFee)
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

func ParseInitializePermissionedPool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseInitializePermissionlessPool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseInitializePermissionlessPoolWithFeeTier(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseEnableOrDisablePool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

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

func ParseRemoveLiquiditySingleSide(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseAddImbalanceLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

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

func ParseSetPoolFees(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseOverrideCurveParam(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseGetPoolInfo(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseBootstrapLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseCreateMintMetadata(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseCreateLockEscrow(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseLock(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseClaimFee(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseCreateConfig(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseCloseConfig(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseInitializePermissionlessConstantProductPoolWithConfig(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseInitializePermissionlessConstantProductPoolWithConfig2(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseInitializeCustomizablePermissionlessConstantProductPool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseUpdateActivationPoint(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseWithdrawProtocolFees(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParseSetWhitelistedVault(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

func ParsePartnerClaimFee(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Default
func ParseDefault(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
