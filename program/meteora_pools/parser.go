package meteora_pools

import (
	"errors"

	"github.com/gagliardetto/solana-go/programs/meteora_pools"
	"github.com/solana-parser/log"
	"github.com/solana-parser/program"
	"github.com/solana-parser/types"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(meteora_pools.ProgramID, meteora_pools.ProgramName, program.Swap, 1, ProgramParser)
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

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	inst, err := meteora_pools.DecodeInstruction(in.RawInstruction.AccountValues, in.RawInstruction.DataBytes)
	if err != nil {
		return err
	}
	id := uint64(inst.TypeID.Uint32())
	parser, ok := Parsers[id]
	if !ok {
		return errors.New("parser not found")
	}
	return parser(inst, in, meta)
}

func ParseInitializePermissionedPool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	log.Logger.Info("ignore parse initialize permissioned pool", "program", meteora_pools.ProgramName)
	return nil
}

func ParseInitializePermissionlessPool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	log.Logger.Info("ignore parse initialize permissionless pool", "program", meteora_pools.ProgramName)
	return nil
}

func ParseInitializePermissionlessPoolWithFeeTier(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	log.Logger.Info("ignore parse initialize permissionless pool with feeTier", "program", meteora_pools.ProgramName)
	return nil
}

func ParseEnableOrDisablePool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseSwap(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*meteora_pools.Swap)
	swap := &types.Swap{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	swap.InputTransfer = in.FindNextTransferByFrom(inst1.GetUserSourceTokenAccount().PublicKey)
	swap.OutputTransfer = in.FindNextTransferByTo(inst1.GetUserDestinationTokenAccount().PublicKey)
	in.Event = []interface{}{swap}
	return nil
}

func ParseRemoveLiquiditySingleSide(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	log.Logger.Info("ignore parse remove liquidity single side", "program", meteora_pools.ProgramName)
	return nil
}

func ParseAddImbalanceLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	// log.Logger.Info("ignore parse add imbalance liquidity", "program", meteora_pools.ProgramName)
	inst1 := inst.Impl.(*meteora_pools.AddImbalanceLiquidity)
	addLiquidity := &types.AddLiquidity{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	addLiquidity.TokenATransfer = in.FindNextTransferByTo(inst1.GetATokenVaultAccount().PublicKey)
	addLiquidity.TokenBTransfer = in.FindNextTransferByTo(inst1.GetBTokenVaultAccount().PublicKey)
	in.Event = []interface{}{addLiquidity}
	return nil
}

func ParseRemoveBalanceLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*meteora_pools.RemoveBalanceLiquidity)
	removeLiquidity := &types.RemoveLiquidity{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	removeLiquidity.TokenATransfer = in.FindNextTransferByFrom(inst1.GetATokenVaultAccount().PublicKey)
	removeLiquidity.TokenBTransfer = in.FindNextTransferByFrom(inst1.GetBTokenVaultAccount().PublicKey)
	in.Event = []interface{}{removeLiquidity}
	return nil
}

func ParseAddBalanceLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*meteora_pools.AddBalanceLiquidity)
	addLiquidity := &types.AddLiquidity{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	addLiquidity.TokenATransfer = in.FindNextTransferByTo(inst1.GetATokenVaultAccount().PublicKey)
	addLiquidity.TokenBTransfer = in.FindNextTransferByTo(inst1.GetBTokenVaultAccount().PublicKey)
	in.Event = []interface{}{addLiquidity}
	return nil
}

func ParseSetPoolFees(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseOverrideCurveParam(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseGetPoolInfo(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseBootstrapLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	log.Logger.Info("ignore parse bootstrap liquidity", "program", meteora_pools.ProgramName)
	return nil
}

func ParseCreateMintMetadata(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseCreateLockEscrow(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseLock(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseClaimFee(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseCreateConfig(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseCloseConfig(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseInitializePermissionlessConstantProductPoolWithConfig(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	log.Logger.Info("ignore parse initialize permissionless constant_product pool with config", "program", meteora_pools.ProgramName)
	return nil
}

func ParseInitializePermissionlessConstantProductPoolWithConfig2(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	// todo, add liquidity
	// find two deposit
	inst1 := inst.Impl.(*meteora_pools.InitializePermissionlessConstantProductPoolWithConfig2)
	addLiquidity := &types.AddLiquidity{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetPayerAccount().PublicKey,
	}
	addLiquidity.TokenATransfer = in.FindNextTransferByTo(inst1.GetATokenVaultAccount().PublicKey)
	addLiquidity.TokenBTransfer = in.FindNextTransferByTo(inst1.GetBTokenVaultAccount().PublicKey)
	in.Event = []interface{}{addLiquidity}
	return nil
}

func ParseInitializeCustomizablePermissionlessConstantProductPool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	log.Logger.Info("ignore parse initialize customizable permissionless constant_product pool", "program", meteora_pools.ProgramName)
	return nil
}

func ParseUpdateActivationPoint(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseWithdrawProtocolFees(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParseSetWhitelistedVault(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

func ParsePartnerClaimFee(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

// Default
func ParseDefault(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

// Fault
func ParseFault(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) error {
	panic("not supported")
}
