package raydium_cp

import (
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/raydium_cp"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(raydium_cp.ProgramID, ProgramParser)
	RegisterParser(uint64(raydium_cp.Instruction_CreateAmmConfig.Uint32()), ParseCreateAmmConfig)
	RegisterParser(uint64(raydium_cp.Instruction_UpdateAmmConfig.Uint32()), ParseUpdateAmmConfig)
	RegisterParser(uint64(raydium_cp.Instruction_UpdatePoolStatus.Uint32()), ParseUpdatePoolStatus)
	RegisterParser(uint64(raydium_cp.Instruction_CollectProtocolFee.Uint32()), ParseCollectProtocolFee)
	RegisterParser(uint64(raydium_cp.Instruction_CollectFundFee.Uint32()), ParseCollectFundFee)
	RegisterParser(uint64(raydium_cp.Instruction_Initialize.Uint32()), ParseInitialize)
	RegisterParser(uint64(raydium_cp.Instruction_Deposit.Uint32()), ParseDeposit)
	RegisterParser(uint64(raydium_cp.Instruction_Withdraw.Uint32()), ParseWithdraw)
	RegisterParser(uint64(raydium_cp.Instruction_SwapBaseInput.Uint32()), ParseSwapBaseInput)
	RegisterParser(uint64(raydium_cp.Instruction_SwapBaseOutput.Uint32()), ParseSwapBaseOutput)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := raydium_cp.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

func ParseCreateAmmConfig(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdateAmmConfig(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdatePoolStatus(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCollectProtocolFee(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCollectFundFee(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitialize(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize", "program", raydium_cp.ProgramName)
}
func ParseDeposit(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_cp.Deposit)
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetPoolStateAccount().PublicKey,
		User:           inst1.Get(0).PublicKey,
		TokenATransfer: in.Children[0].Event[0].(*types.Transfer),
		TokenBTransfer: in.Children[1].Event[0].(*types.Transfer),
		TokenLpMint:    in.Children[2].Event[0].(*types.MintTo),
	}
	in.Event = []interface{}{addLiquidity}
}
func ParseWithdraw(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_cp.Withdraw)
	// child 1 : transfer
	// child 2 : transfer
	removeLiquidity := &types.RemoveLiquidity{
		Pool:           inst1.GetPoolStateAccount().PublicKey,
		User:           inst1.GetAuthorityAccount().PublicKey,
		TokenLpBurn:    in.Children[0].Event[0].(*types.Burn),
		TokenATransfer: in.Children[1].Event[0].(*types.Transfer),
		TokenBTransfer: in.Children[2].Event[0].(*types.Transfer),
	}
	in.Event = []interface{}{removeLiquidity}
}
func ParseSwapBaseInput(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_cp.SwapBaseInput)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	user := inst1.GetInputTokenAccountAccount().PublicKey
	if owner, ok := meta.TokenOwner[user]; ok {
		user = owner
	}
	swap := &types.Swap{
		Pool:           inst1.GetPoolStateAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
		User:           user,
	}
	in.Event = []interface{}{swap}
}
func ParseSwapBaseOutput(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_cp.SwapBaseOutput)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	user := inst1.GetInputTokenAccountAccount().PublicKey
	if owner, ok := meta.TokenOwner[user]; ok {
		user = owner
	}
	swap := &types.Swap{
		Pool:           inst1.GetPoolStateAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
		User:           user,
	}
	in.Event = []interface{}{swap}
}

// Default
func ParseDefault(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
