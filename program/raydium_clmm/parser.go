package raydium_clmm

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
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
	program.RegisterParser(raydium_clmm.ProgramID, raydium_clmm.ProgramName, ProgramParser)
	RegisterParser(uint64(raydium_clmm.Instruction_CreateAmmConfig.Uint32()), ParseCreateAmmConfig)
	RegisterParser(uint64(raydium_clmm.Instruction_UpdateAmmConfig.Uint32()), ParseUpdateAmmConfig)
	RegisterParser(uint64(raydium_clmm.Instruction_CreatePool.Uint32()), ParseCreatePool)
	RegisterParser(uint64(raydium_clmm.Instruction_UpdatePoolStatus.Uint32()), ParseUpdatePoolStatus)
	RegisterParser(uint64(raydium_clmm.Instruction_CreateOperationAccount.Uint32()), ParseCreateOperationAccount)
	RegisterParser(uint64(raydium_clmm.Instruction_UpdateOperationAccount.Uint32()), ParseUpdateOperationAccount)
	RegisterParser(uint64(raydium_clmm.Instruction_TransferRewardOwner.Uint32()), ParseTransferRewardOwner)
	RegisterParser(uint64(raydium_clmm.Instruction_InitializeReward.Uint32()), ParseInitializeReward)
	RegisterParser(uint64(raydium_clmm.Instruction_CollectRemainingRewards.Uint32()), ParseCollectRemainingRewards)
	RegisterParser(uint64(raydium_clmm.Instruction_UpdateRewardInfos.Uint32()), ParseUpdateRewardInfos)
	RegisterParser(uint64(raydium_clmm.Instruction_SetRewardParams.Uint32()), ParseSetRewardParams)
	RegisterParser(uint64(raydium_clmm.Instruction_CollectProtocolFee.Uint32()), ParseCollectProtocolFee)
	RegisterParser(uint64(raydium_clmm.Instruction_CollectFundFee.Uint32()), ParseCollectFundFee)
	RegisterParser(uint64(raydium_clmm.Instruction_OpenPosition.Uint32()), ParseOpenPosition)
	RegisterParser(uint64(raydium_clmm.Instruction_OpenPositionV2.Uint32()), ParseOpenPositionV2)
	RegisterParser(uint64(raydium_clmm.Instruction_OpenPositionWithToken22Nft.Uint32()), ParseOpenPositionWithToken22Nft)
	RegisterParser(uint64(raydium_clmm.Instruction_ClosePosition.Uint32()), ParseClosePosition)
	RegisterParser(uint64(raydium_clmm.Instruction_IncreaseLiquidity.Uint32()), ParseIncreaseLiquidity)
	RegisterParser(uint64(raydium_clmm.Instruction_IncreaseLiquidityV2.Uint32()), ParseIncreaseLiquidityV2)
	RegisterParser(uint64(raydium_clmm.Instruction_DecreaseLiquidity.Uint32()), ParseDecreaseLiquidity)
	RegisterParser(uint64(raydium_clmm.Instruction_DecreaseLiquidityV2.Uint32()), ParseDecreaseLiquidityV2)
	RegisterParser(uint64(raydium_clmm.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(raydium_clmm.Instruction_SwapV2.Uint32()), ParseSwapV2)
	RegisterParser(uint64(raydium_clmm.Instruction_SwapRouterBaseIn.Uint32()), ParseSwapRouterBaseIn)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	inst, err := raydium_clmm.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

func ParseCreateAmmConfig(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdateAmmConfig(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCreatePool(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_clmm.CreatePool)
	createPool := &types.CreatePool{
		Dex:     in.Instruction.ProgramId,
		Pool:    inst1.GetPoolStateAccount().PublicKey,
		User:    inst1.GetPoolCreatorAccount().PublicKey,
		TokenA:  inst1.GetTokenMint0Account().PublicKey,
		TokenB:  inst1.GetTokenMint1Account().PublicKey,
		TokenLP: solana.PublicKey{},
		VaultA:  inst1.GetTokenVault0Account().PublicKey,
		VaultB:  inst1.GetTokenVault1Account().PublicKey,
		VaultLP: solana.PublicKey{},
	}
	in.Event = []interface{}{createPool}
}
func ParseUpdatePoolStatus(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCreateOperationAccount(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdateOperationAccount(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseTransferRewardOwner(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitializeReward(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCollectRemainingRewards(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdateRewardInfos(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetRewardParams(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCollectProtocolFee(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCollectFundFee(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseOpenPosition(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse open position", "program", raydium_clmm.ProgramName)
}
func ParseOpenPositionV2(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse open position v2", "program", raydium_clmm.ProgramName)
}
func ParseOpenPositionWithToken22Nft(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_clmm.OpenPositionWithToken22Nft)
	transfers := in.FindChildrenTransfers()
	addLiquidity := &types.AddLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolStateAccount().PublicKey,
		User: inst1.GetPayerAccount().PublicKey,
	}
	if len(transfers) >= 1 {
		addLiquidity.TokenATransfer = transfers[0]
	}
	if len(transfers) >= 2 {
		addLiquidity.TokenBTransfer = transfers[1]
	}
	in.Event = []interface{}{addLiquidity}
}
func ParseClosePosition(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	// close all accounts
}
func ParseIncreaseLiquidity(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse increase liquidity", "program", raydium_clmm.ProgramName)
}
func ParseIncreaseLiquidityV2(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_clmm.IncreaseLiquidityV2)
	transfers := in.FindChildrenTransfers()
	addLiquidity := &types.AddLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolStateAccount().PublicKey,
		User: inst1.GetNftOwnerAccount().PublicKey,
	}
	if len(transfers) >= 1 {
		addLiquidity.TokenATransfer = transfers[0]
	}
	if len(transfers) >= 2 {
		addLiquidity.TokenBTransfer = transfers[1]
	}
	in.Event = []interface{}{addLiquidity}
}
func ParseDecreaseLiquidity(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse decrease liquidity", "program", raydium_clmm.ProgramName)
}
func ParseDecreaseLiquidityV2(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_clmm.DecreaseLiquidityV2)
	removeLiquidity := &types.RemoveLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolStateAccount().PublicKey,
		User: inst1.GetNftOwnerAccount().PublicKey,
	}
	if len(in.Children) >= 1 {
		removeLiquidity.TokenATransfer = in.Children[0].Event[0].(*types.Transfer)
	}
	if len(in.Children) >= 2 {
		removeLiquidity.TokenBTransfer = in.Children[1].Event[0].(*types.Transfer)
	}
	in.Event = []interface{}{removeLiquidity}
}
func ParseSwap(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_clmm.Swap)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	//
	swap := &types.Swap{
		Dex:            in.Instruction.ProgramId,
		Pool:           inst1.GetPoolStateAccount().PublicKey,
		User:           inst1.GetPayerAccount().PublicKey,
		InputTransfer:  t1,
		OutputTransfer: t2,
	}
	in.Event = []interface{}{swap}
}
func ParseSwapV2(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_clmm.SwapV2)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	swap := &types.Swap{
		Dex:            in.Instruction.ProgramId,
		Pool:           inst1.GetPoolStateAccount().PublicKey,
		User:           inst1.GetPayerAccount().PublicKey,
		InputTransfer:  t1,
		OutputTransfer: t2,
	}
	in.Event = []interface{}{swap}
}
func ParseSwapRouterBaseIn(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse swap router base in", "program", raydium_clmm.ProgramName)
}

// Default
func ParseDefault(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *raydium_clmm.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
