package whirlpool

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/whirlpool"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

var (
	Instruction_OpenPositionWithTokenExtensions  = ag_binary.TypeID([8]byte{212, 47, 95, 92, 114, 102, 131, 250})
	Instruction_ClosePositionWithTokenExtensions = ag_binary.TypeID([8]byte{1, 182, 135, 59, 155, 25, 99, 223})
)

func init() {
	program.RegisterParser(whirlpool.ProgramID, ProgramParser)
	RegisterParser(uint64(whirlpool.Instruction_InitializeConfig.Uint32()), ParseInitializeConfig)
	RegisterParser(uint64(whirlpool.Instruction_InitializePool.Uint32()), ParseInitializePool)
	RegisterParser(uint64(whirlpool.Instruction_InitializeTickArray.Uint32()), ParseInitializeTickArray)
	RegisterParser(uint64(whirlpool.Instruction_InitializeFeeTier.Uint32()), ParseInitializeFeeTier)
	RegisterParser(uint64(whirlpool.Instruction_InitializeReward.Uint32()), ParseInitializeReward)
	RegisterParser(uint64(whirlpool.Instruction_SetRewardEmissions.Uint32()), ParseSetRewardEmissions)
	RegisterParser(uint64(whirlpool.Instruction_OpenPosition.Uint32()), ParseOpenPosition)
	RegisterParser(uint64(whirlpool.Instruction_OpenPositionWithMetadata.Uint32()), ParseOpenPositionWithMetadata)
	RegisterParser(uint64(whirlpool.Instruction_IncreaseLiquidity.Uint32()), ParseIncreaseLiquidity)
	RegisterParser(uint64(whirlpool.Instruction_DecreaseLiquidity.Uint32()), ParseDecreaseLiquidity)
	RegisterParser(uint64(whirlpool.Instruction_UpdateFeesAndRewards.Uint32()), ParseUpdateFeesAndRewards)
	RegisterParser(uint64(whirlpool.Instruction_CollectFees.Uint32()), ParseCollectFees)
	RegisterParser(uint64(whirlpool.Instruction_CollectReward.Uint32()), ParseCollectReward)
	RegisterParser(uint64(whirlpool.Instruction_CollectProtocolFees.Uint32()), ParseCollectProtocolFees)
	RegisterParser(uint64(whirlpool.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(whirlpool.Instruction_ClosePosition.Uint32()), ParseClosePosition)
	RegisterParser(uint64(whirlpool.Instruction_SetDefaultFeeRate.Uint32()), ParseSetDefaultFeeRate)
	RegisterParser(uint64(whirlpool.Instruction_SetDefaultProtocolFeeRate.Uint32()), ParseSetDefaultProtocolFeeRate)
	RegisterParser(uint64(whirlpool.Instruction_SetFeeRate.Uint32()), ParseSetFeeRate)
	RegisterParser(uint64(whirlpool.Instruction_SetProtocolFeeRate.Uint32()), ParseSetProtocolFeeRate)
	RegisterParser(uint64(whirlpool.Instruction_SetFeeAuthority.Uint32()), ParseSetFeeAuthority)
	RegisterParser(uint64(whirlpool.Instruction_SetCollectProtocolFeesAuthority.Uint32()), ParseSetCollectProtocolFeesAuthority)
	RegisterParser(uint64(whirlpool.Instruction_SetRewardAuthority.Uint32()), ParseSetRewardAuthority)
	RegisterParser(uint64(whirlpool.Instruction_SetRewardAuthorityBySuperAuthority.Uint32()), ParseSetRewardAuthorityBySuperAuthority)
	RegisterParser(uint64(whirlpool.Instruction_SetRewardEmissionsSuperAuthority.Uint32()), ParseSetRewardEmissionsSuperAuthority)
	RegisterParser(uint64(whirlpool.Instruction_TwoHopSwap.Uint32()), ParseTwoHopSwap)
	RegisterParser(uint64(whirlpool.Instruction_InitializePositionBundle.Uint32()), ParseInitializePositionBundle)
	RegisterParser(uint64(whirlpool.Instruction_InitializePositionBundleWithMetadata.Uint32()), ParseInitializePositionBundleWithMetadata)
	RegisterParser(uint64(whirlpool.Instruction_DeletePositionBundle.Uint32()), ParseDeletePositionBundle)
	RegisterParser(uint64(whirlpool.Instruction_OpenBundledPosition.Uint32()), ParseOpenBundledPosition)
	RegisterParser(uint64(whirlpool.Instruction_CloseBundledPosition.Uint32()), ParseCloseBundledPosition)
	RegisterParser(uint64(whirlpool.Instruction_CollectFeesV2.Uint32()), ParseCollectFeesV2)
	RegisterParser(uint64(whirlpool.Instruction_CollectProtocolFeesV2.Uint32()), ParseCollectProtocolFeesV2)
	RegisterParser(uint64(whirlpool.Instruction_CollectRewardV2.Uint32()), ParseCollectRewardV2)
	RegisterParser(uint64(whirlpool.Instruction_DecreaseLiquidityV2.Uint32()), ParseDecreaseLiquidityV2)
	RegisterParser(uint64(whirlpool.Instruction_IncreaseLiquidityV2.Uint32()), ParseIncreaseLiquidityV2)
	RegisterParser(uint64(whirlpool.Instruction_InitializePoolV2.Uint32()), ParseInitializePoolV2)
	RegisterParser(uint64(whirlpool.Instruction_InitializeRewardV2.Uint32()), ParseInitializeRewardV2)
	RegisterParser(uint64(whirlpool.Instruction_SetRewardEmissionsV2.Uint32()), ParseSetRewardEmissionsV2)
	RegisterParser(uint64(whirlpool.Instruction_SwapV2.Uint32()), ParseSwapV2)
	RegisterParser(uint64(whirlpool.Instruction_TwoHopSwapV2.Uint32()), ParseTwoHopSwapV2)
	RegisterParser(uint64(whirlpool.Instruction_InitializeConfigExtension.Uint32()), ParseInitializeConfigExtension)
	RegisterParser(uint64(whirlpool.Instruction_SetConfigExtensionAuthority.Uint32()), ParseSetConfigExtensionAuthority)
	RegisterParser(uint64(whirlpool.Instruction_SetTokenBadgeAuthority.Uint32()), ParseSetTokenBadgeAuthority)
	RegisterParser(uint64(whirlpool.Instruction_InitializeTokenBadge.Uint32()), ParseInitializeTokenBadge)
	RegisterParser(uint64(whirlpool.Instruction_DeleteTokenBadge.Uint32()), ParseDeleteTokenBadge)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	dec := ag_binary.NewBorshDecoder(in.Instruction.Data)
	typeID, err := dec.ReadTypeID()
	if typeID == Instruction_ClosePositionWithTokenExtensions || typeID == Instruction_OpenPositionWithTokenExtensions {
		return nil
	}
	inst, err := whirlpool.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

func ParseInitializeConfig(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitializePool(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	// log.Logger.Info("ignore parse initialize pool", "program", whirlpool.ProgramName)
	inst1 := inst.Impl.(*whirlpool.InitializePool)
	createPool := &types.CreatePool{
		Pool:    inst1.GetWhirlpoolAccount().PublicKey,
		User:    inst1.GetFunderAccount().PublicKey,
		TokenA:  inst1.GetTokenMintAAccount().PublicKey,
		TokenB:  inst1.GetTokenMintBAccount().PublicKey,
		TokenLP: solana.PublicKey{},
		VaultA:  inst1.GetTokenVaultAAccount().PublicKey,
		VaultB:  inst1.GetTokenVaultBAccount().PublicKey,
		VaultLP: solana.PublicKey{},
	}
	in.Event = []interface{}{createPool}
}
func ParseInitializeTickArray(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitializeFeeTier(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitializeReward(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetRewardEmissions(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseOpenPosition(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	// only create all accounts
	// log.Logger.Info("ignore parse open position", "program", whirlpool.ProgramName)
}
func ParseOpenPositionWithMetadata(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	// log.Logger.Info("ignore parse open position with metadata", "program", whirlpool.ProgramName)
}
func ParseIncreaseLiquidity(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*whirlpool.IncreaseLiquidity)
	// child 1 : transfer
	// child 2 : transfer
	transfers := in.FindChildrenTransfers()
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetWhirlpoolAccount().PublicKey,
		User:           inst1.GetPositionAuthorityAccount().PublicKey,
		TokenATransfer: transfers[0],
		TokenBTransfer: transfers[1],
	}
	in.Event = []interface{}{addLiquidity}

}
func ParseDecreaseLiquidity(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*whirlpool.DecreaseLiquidity)
	// child 1 : transfer
	// child 2 : transfer
	transfers := in.FindChildrenTransfers()
	removeLiquidity := &types.RemoveLiquidity{
		Pool:           inst1.GetWhirlpoolAccount().PublicKey,
		User:           inst1.GetPositionAuthorityAccount().PublicKey,
		TokenATransfer: transfers[0],
		TokenBTransfer: transfers[1],
	}
	in.Event = []interface{}{removeLiquidity}
}
func ParseUpdateFeesAndRewards(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCollectFees(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCollectReward(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCollectProtocolFees(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSwap(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*whirlpool.Swap)
	// child 1 : input transfer
	// child 2 : output transfer
	transfers := in.FindChildrenTransfers()
	swap := &types.Swap{
		Pool:           inst1.GetWhirlpoolAccount().PublicKey,
		User:           inst1.GetTokenAuthorityAccount().PublicKey,
		InputTransfer:  transfers[0],
		OutputTransfer: transfers[1],
	}
	in.Event = []interface{}{swap}
}
func ParseClosePosition(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	// close all accounts
}
func ParseSetDefaultFeeRate(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetDefaultProtocolFeeRate(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetFeeRate(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetProtocolFeeRate(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetFeeAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetCollectProtocolFeesAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetRewardAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetRewardAuthorityBySuperAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetRewardEmissionsSuperAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseTwoHopSwap(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*whirlpool.TwoHopSwap)
	// child 1 : transfer
	// child 2 : transfer
	// child 3 : transfer
	transfers := in.FindChildrenTransfers()
	swap := &types.Swap{
		Pool:           inst1.GetWhirlpoolOneAccount().PublicKey,
		User:           inst1.GetTokenAuthorityAccount().PublicKey,
		InputTransfer:  transfers[0],
		OutputTransfer: transfers[2],
	}
	in.Event = []interface{}{swap}
}
func ParseInitializePositionBundle(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	//log.Logger.Info("ignore parse initialize position bundle", "program", whirlpool.ProgramName)
}
func ParseInitializePositionBundleWithMetadata(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	//log.Logger.Info("ignore parse initialize position bundle", "program", whirlpool.ProgramName)
}
func ParseDeletePositionBundle(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	//log.Logger.Info("ignore parse delete position bundle", "program", whirlpool.ProgramName)
}
func ParseOpenBundledPosition(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	//log.Logger.Info("ignore parse open bundled position", "program", whirlpool.ProgramName)
}
func ParseCloseBundledPosition(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	//log.Logger.Info("ignore parse close bundle position", "program", whirlpool.ProgramName)
}
func ParseCollectFeesV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCollectProtocolFeesV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCollectRewardV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseDecreaseLiquidityV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*whirlpool.DecreaseLiquidityV2)
	// child 1 : transfer
	// child 2 : transfer
	transfers := in.FindChildrenTransfers()
	removeLiquidity := &types.RemoveLiquidity{
		Pool:           inst1.GetWhirlpoolAccount().PublicKey,
		User:           inst1.GetPositionAuthorityAccount().PublicKey,
		TokenATransfer: transfers[0],
		TokenBTransfer: transfers[1],
	}
	in.Event = []interface{}{removeLiquidity}
}
func ParseIncreaseLiquidityV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*whirlpool.IncreaseLiquidityV2)
	// child 1 : transfer
	// child 2 : transfer
	transfers := in.FindChildrenTransfers()
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetWhirlpoolAccount().PublicKey,
		User:           inst1.GetPositionAuthorityAccount().PublicKey,
		TokenATransfer: transfers[0],
		TokenBTransfer: transfers[1],
	}
	in.Event = []interface{}{addLiquidity}
}
func ParseInitializePoolV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	// log.Logger.Info("ignore parse initialize pool v2", "program", whirlpool.ProgramName)
	inst1 := inst.Impl.(*whirlpool.InitializePoolV2)
	createPool := &types.CreatePool{
		Pool:    inst1.GetWhirlpoolAccount().PublicKey,
		User:    inst1.GetFunderAccount().PublicKey,
		TokenA:  inst1.GetTokenMintAAccount().PublicKey,
		TokenB:  inst1.GetTokenMintBAccount().PublicKey,
		TokenLP: solana.PublicKey{},
		VaultA:  inst1.GetTokenVaultAAccount().PublicKey,
		VaultB:  inst1.GetTokenVaultBAccount().PublicKey,
		VaultLP: solana.PublicKey{},
	}
	in.Event = []interface{}{createPool}
}
func ParseInitializeRewardV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	//log.Logger.Info("ignore parse initialize reward v2", "program", whirlpool.ProgramName)
}
func ParseSetRewardEmissionsV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	//log.Logger.Info("ignore parse set reward emissions v2", "program", whirlpool.ProgramName)
}
func ParseSwapV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*whirlpool.SwapV2)
	// child 1 : input transfer
	// child 2 : output transfer
	transfers := in.FindChildrenTransfers()
	swap := &types.Swap{
		Pool:           inst1.GetWhirlpoolAccount().PublicKey,
		User:           inst1.GetTokenAuthorityAccount().PublicKey,
		InputTransfer:  transfers[0],
		OutputTransfer: transfers[1],
	}
	in.Event = []interface{}{swap}
}
func ParseTwoHopSwapV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*whirlpool.TwoHopSwapV2)
	// child 1 : transfer
	// child 2 : transfer
	// child 3 : transfer
	transfers := in.FindChildrenTransfers()
	swap := &types.Swap{
		Pool:           inst1.GetWhirlpoolOneAccount().PublicKey,
		User:           inst1.GetTokenAuthorityAccount().PublicKey,
		InputTransfer:  transfers[0],
		OutputTransfer: transfers[2],
	}
	in.Event = []interface{}{swap}
}
func ParseInitializeConfigExtension(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetConfigExtensionAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetTokenBadgeAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitializeTokenBadge(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseDeleteTokenBadge(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
}

// Default
func ParseDefault(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
