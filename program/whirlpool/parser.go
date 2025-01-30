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

func ParseInitializeConfig(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializePool(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializeTickArray(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializeFeeTier(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializeReward(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetRewardEmissions(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseOpenPosition(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseOpenPositionWithMetadata(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseIncreaseLiquidity(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	// child 1 : transfer
	// child 2 : transfer
	transfer1 := in.Children[0].Event[0]
	transfer2 := in.Children[0].Event[0]
	in.Event = []interface{}{transfer1, transfer2}
}
func ParseDecreaseLiquidity(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	// child 1 : transfer
	// child 2 : transfer
	transfer1 := in.Children[0].Event[0]
	transfer2 := in.Children[0].Event[0]
	in.Event = []interface{}{transfer1, transfer2}
}
func ParseUpdateFeesAndRewards(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCollectFees(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCollectReward(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCollectProtocolFees(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
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
func ParseClosePosition(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetDefaultFeeRate(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetDefaultProtocolFeeRate(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetFeeRate(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetProtocolFeeRate(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetFeeAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetCollectProtocolFeesAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetRewardAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetRewardAuthorityBySuperAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetRewardEmissionsSuperAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseTwoHopSwap(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializePositionBundle(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializePositionBundleWithMetadata(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseDeletePositionBundle(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseOpenBundledPosition(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCloseBundledPosition(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCollectFeesV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCollectProtocolFeesV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCollectRewardV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseDecreaseLiquidityV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseIncreaseLiquidityV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializePoolV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializeRewardV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetRewardEmissionsV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
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
func ParseTwoHopSwapV2(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializeConfigExtension(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetConfigExtensionAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetTokenBadgeAuthority(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializeTokenBadge(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseDeleteTokenBadge(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Default
func ParseDefault(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *whirlpool.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
