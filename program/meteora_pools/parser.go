package meteora_pools

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/log"
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
	program.RegisterParser(meteora_pools.ProgramID, meteora_pools.ProgramName, program.Swap, ProgramParser)
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
	inst, err := meteora_pools.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
	if err != nil {
		return err
	}
	id := uint64(inst.TypeID.Uint32())
	parser, ok := Parsers[id]
	if !ok {
		return errors.New("parser not found")
	}
	parser(inst, in, meta)
	return nil
}

func ParseInitializePermissionedPool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize permissioned pool", "program", meteora_pools.ProgramName)
}

func ParseInitializePermissionlessPool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize permissionless pool", "program", meteora_pools.ProgramName)
}

func ParseInitializePermissionlessPoolWithFeeTier(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize permissionless pool with feeTier", "program", meteora_pools.ProgramName)
}

func ParseEnableOrDisablePool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseSwap(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_pools.Swap)
	swap := &types.Swap{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	if *inst1.InAmount != 0 {
		// the transfer is execute by vault deposit & withdraw & this first transfer is fee
		transfers := in.FindChildrenTransfers()
		for _, transfer := range transfers {
			if transfer.To == inst1.GetATokenVaultAccount().PublicKey || transfer.To == inst1.GetBTokenVaultAccount().PublicKey {
				swap.InputTransfer = transfer
			}
			if transfer.From == inst1.GetATokenVaultAccount().PublicKey || transfer.From == inst1.GetBTokenVaultAccount().PublicKey {
				swap.OutputTransfer = transfer
			}
		}
	}
	in.Event = []interface{}{swap}
}

func ParseRemoveLiquiditySingleSide(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse remove liquidity single side", "program", meteora_pools.ProgramName)
}

func ParseAddImbalanceLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	// log.Logger.Info("ignore parse add imbalance liquidity", "program", meteora_pools.ProgramName)
	inst1 := inst.Impl.(*meteora_pools.AddImbalanceLiquidity)
	addLiquidity := &types.AddLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.To == inst1.GetATokenVaultAccount().PublicKey {
			addLiquidity.TokenATransfer = transfer
		}
		if transfer.To == inst1.GetBTokenVaultAccount().PublicKey {
			addLiquidity.TokenBTransfer = transfer
		}
	}
	in.Event = []interface{}{addLiquidity}
}

func ParseRemoveBalanceLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_pools.RemoveBalanceLiquidity)
	removeLiquidity := &types.RemoveLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.From == inst1.GetATokenVaultAccount().PublicKey {
			removeLiquidity.TokenATransfer = transfer
		}
		if transfer.From == inst1.GetBTokenVaultAccount().PublicKey {
			removeLiquidity.TokenBTransfer = transfer
		}
	}
	in.Event = []interface{}{removeLiquidity}
}

func ParseAddBalanceLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*meteora_pools.AddBalanceLiquidity)
	addLiquidity := &types.AddLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.To == inst1.GetATokenVaultAccount().PublicKey {
			addLiquidity.TokenATransfer = transfer
		}
		if transfer.To == inst1.GetBTokenVaultAccount().PublicKey {
			addLiquidity.TokenBTransfer = transfer
		}
	}
	in.Event = []interface{}{addLiquidity}
}

func ParseSetPoolFees(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseOverrideCurveParam(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseGetPoolInfo(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseBootstrapLiquidity(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse bootstrap liquidity", "program", meteora_pools.ProgramName)
}

func ParseCreateMintMetadata(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseCreateLockEscrow(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseLock(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseClaimFee(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseCreateConfig(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseCloseConfig(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseInitializePermissionlessConstantProductPoolWithConfig(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize permissionless constant_product pool with config", "program", meteora_pools.ProgramName)
}

func ParseInitializePermissionlessConstantProductPoolWithConfig2(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	// todo, add liquidity
	// find two deposit
	inst1 := inst.Impl.(*meteora_pools.InitializePermissionlessConstantProductPoolWithConfig2)
	addLiquidity := &types.AddLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetPayerAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.To == inst1.GetATokenVaultAccount().PublicKey {
			addLiquidity.TokenATransfer = transfer
		}
		if transfer.To == inst1.GetBTokenVaultAccount().PublicKey {
			addLiquidity.TokenBTransfer = transfer
		}
	}
	in.Event = []interface{}{addLiquidity}
}

func ParseInitializeCustomizablePermissionlessConstantProductPool(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize customizable permissionless constant_product pool", "program", meteora_pools.ProgramName)
}

func ParseUpdateActivationPoint(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseWithdrawProtocolFees(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParseSetWhitelistedVault(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

func ParsePartnerClaimFee(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
}

// Default
func ParseDefault(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *meteora_pools.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
