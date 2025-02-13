package raydium_cp

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
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
	program.RegisterParser(raydium_cp.ProgramID, raydium_cp.ProgramName, program.Swap, ProgramParser)
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

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	inst, err := raydium_cp.DecodeInstruction(in.AccountMetas(meta.Accounts), in.Instruction.Data)
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
	// log.Logger.Info("ignore parse initialize", "program", raydium_cp.ProgramName)
	inst1 := inst.Impl.(*raydium_cp.Initialize)
	createPool := &types.CreatePool{
		Dex:     in.Instruction.ProgramId,
		Pool:    inst1.GetPoolStateAccount().PublicKey,
		User:    inst1.GetCreatorAccount().PublicKey,
		TokenA:  inst1.GetToken0MintAccount().PublicKey,
		TokenB:  inst1.GetToken1MintAccount().PublicKey,
		TokenLP: inst1.GetLpMintAccount().PublicKey,
		VaultA:  inst1.GetToken0VaultAccount().PublicKey,
		VaultB:  inst1.GetToken1VaultAccount().PublicKey,
		VaultLP: solana.PublicKey{},
	}
	addLiquidity := &types.AddLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolStateAccount().PublicKey,
		User: inst1.GetCreatorAccount().PublicKey,
	}
	// the latest three transfer
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.To == inst1.GetToken0VaultAccount().PublicKey {
			addLiquidity.TokenATransfer = transfer
		}
		if transfer.To == inst1.GetToken1VaultAccount().PublicKey {
			addLiquidity.TokenBTransfer = transfer
		}
	}
	in.Event = []interface{}{createPool, addLiquidity}
}
func ParseDeposit(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_cp.Deposit)
	addLiquidity := &types.AddLiquidity{
		Dex:         in.Instruction.ProgramId,
		Pool:        inst1.GetPoolStateAccount().PublicKey,
		User:        inst1.GetOwnerAccount().PublicKey,
		TokenLpMint: in.Children[2].Event[0].(*types.MintTo),
	}
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.To == inst1.GetToken0VaultAccount().PublicKey {
			addLiquidity.TokenATransfer = transfer
		}
		if transfer.To == inst1.GetToken1VaultAccount().PublicKey {
			addLiquidity.TokenBTransfer = transfer
		}
	}
	in.Event = []interface{}{addLiquidity}
}
func ParseWithdraw(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_cp.Withdraw)
	removeLiquidity := &types.RemoveLiquidity{
		Dex:            in.Instruction.ProgramId,
		Pool:           inst1.GetPoolStateAccount().PublicKey,
		User:           inst1.GetOwnerAccount().PublicKey,
		TokenBTransfer: in.Children[2].Event[0].(*types.Transfer),
	}
	// child 1 : transfer
	// child 2 : transfer
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.From == inst1.GetToken0VaultAccount().PublicKey {
			removeLiquidity.TokenATransfer = transfer
		}
		if transfer.From == inst1.GetToken1VaultAccount().PublicKey {
			removeLiquidity.TokenBTransfer = transfer
		}
	}
	in.Event = []interface{}{removeLiquidity}
}
func ParseSwapBaseInput(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_cp.SwapBaseInput)
	swap := &types.Swap{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolStateAccount().PublicKey,
		User: inst1.GetPayerAccount().PublicKey,
	}
	if *inst1.AmountIn > 0 {
		swap.InputTransfer = in.Children[0].Event[0].(*types.Transfer)
		swap.OutputTransfer = in.Children[1].Event[0].(*types.Transfer)
	}
	in.Event = []interface{}{swap}
}
func ParseSwapBaseOutput(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_cp.SwapBaseOutput)
	swap := &types.Swap{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolStateAccount().PublicKey,
		User: inst1.GetPayerAccount().PublicKey,
	}
	if *inst1.AmountOut > 0 {
		swap.InputTransfer = in.Children[0].Event[0].(*types.Transfer)
		swap.OutputTransfer = in.Children[1].Event[0].(*types.Transfer)
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
