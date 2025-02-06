package stable_swap

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/stable_swap"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(stable_swap.ProgramID, stable_swap.ProgramName, ProgramParser)
	RegisterParser(uint64(stable_swap.Instruction_AcceptOwner.Uint32()), ParseAcceptOwner)
	RegisterParser(uint64(stable_swap.Instruction_ApproveStrategy.Uint32()), ParseApproveStrategy)
	RegisterParser(uint64(stable_swap.Instruction_ChangeAmpFactor.Uint32()), ParseChangeAmpFactor)
	RegisterParser(uint64(stable_swap.Instruction_ChangeSwapFee.Uint32()), ParseChangeSwapFee)
	RegisterParser(uint64(stable_swap.Instruction_CreateStrategy.Uint32()), ParseCreateStrategy)
	RegisterParser(uint64(stable_swap.Instruction_Deposit.Uint32()), ParseDeposit)
	RegisterParser(uint64(stable_swap.Instruction_ExecStrategy.Uint32()), ParseExecStrategy)
	RegisterParser(uint64(stable_swap.Instruction_Initialize.Uint32()), ParseInitialize)
	RegisterParser(uint64(stable_swap.Instruction_Pause.Uint32()), ParsePause)
	RegisterParser(uint64(stable_swap.Instruction_RejectOwner.Uint32()), ParseRejectOwner)
	RegisterParser(uint64(stable_swap.Instruction_Shutdown.Uint32()), ParseShutdown)
	RegisterParser(uint64(stable_swap.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(stable_swap.Instruction_SwapV2.Uint32()), ParseSwapV2)
	RegisterParser(uint64(stable_swap.Instruction_TransferOwner.Uint32()), ParseTransferOwner)
	RegisterParser(uint64(stable_swap.Instruction_Unpause.Uint32()), ParseUnpause)
	RegisterParser(uint64(stable_swap.Instruction_Withdraw.Uint32()), ParseWithdraw)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	inst, err := stable_swap.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

func ParseAcceptOwner(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseApproveStrategy(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseChangeAmpFactor(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseChangeSwapFee(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCreateStrategy(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseDeposit(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*stable_swap.Deposit)
	addLiquidity := &types.AddLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.To == inst1.GetVaultTokenAAccount().PublicKey {
			addLiquidity.TokenATransfer = transfer
		}
	}
	in.Event = []interface{}{addLiquidity}
	// log.Logger.Info("ignore parse deposit", "program", stable_swap.ProgramName)
}
func ParseExecStrategy(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitialize(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize", "program", stable_swap.ProgramName)
}
func ParsePause(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseRejectOwner(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseShutdown(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSwap(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*stable_swap.Swap)
	swap := &types.Swap{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	if *inst1.AmountIn > 0 {
		swap.InputTransfer = in.Children[0].Event[0].(*types.Transfer)
		swap.OutputTransfer = in.Children[1].Event[0].(*types.Transfer)
	}
	in.Event = []interface{}{swap}
}
func ParseSwapV2(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*stable_swap.SwapV2)
	swap := &types.Swap{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	if *inst1.AmountIn > 0 {
		swap.InputTransfer = in.Children[0].Event[0].(*types.Transfer)
		swap.OutputTransfer = in.Children[1].Event[0].(*types.Transfer)
	}
	in.Event = []interface{}{swap}
}
func ParseTransferOwner(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUnpause(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseWithdraw(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*stable_swap.Withdraw)
	removeLiquidity := &types.RemoveLiquidity{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	transfers := in.FindChildrenTransfers()
	for _, transfer := range transfers {
		if transfer.From == inst1.GetVaultTokenAAccount().PublicKey {
			removeLiquidity.TokenATransfer = transfer
		}
	}
	in.Event = []interface{}{removeLiquidity}
	// log.Logger.Info("ignore parse withdraw", "program", stable_swap.ProgramName)
}

// Default
func ParseDefault(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
