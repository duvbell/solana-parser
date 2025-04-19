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

type Parser func(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(stable_swap.ProgramID, stable_swap.ProgramName, program.StableSwap, 1, ProgramParser)
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
	inst, err := stable_swap.DecodeInstruction(in.RawInstruction.AccountValues, in.RawInstruction.DataBytes)
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

func ParseAcceptOwner(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseApproveStrategy(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseChangeAmpFactor(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseChangeSwapFee(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseCreateStrategy(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseDeposit(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*stable_swap.Deposit)
	addLiquidity := &types.AddLiquidity{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	addLiquidity.TokenATransfer = in.FindChildTransferByTo(inst1.GetVaultTokenAAccount().PublicKey)
	in.Event = []interface{}{addLiquidity}
	return nil
}
func ParseExecStrategy(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseInitialize(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	log.Logger.Info("ignore parse initialize", "program", stable_swap.ProgramName)
	return nil
}
func ParsePause(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseRejectOwner(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseShutdown(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseSwap(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*stable_swap.Swap)
	swap := &types.Swap{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	swap.InputTransfer = in.FindChildTransferByTo(inst1.GetUserTokenInAccount().PublicKey)
	swap.OutputTransfer = in.FindChildTransferByFrom(inst1.GetUserTokenOutAccount().PublicKey)
	in.Event = []interface{}{swap}
	return nil
}
func ParseSwapV2(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*stable_swap.SwapV2)
	swap := &types.Swap{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	swap.InputTransfer = in.FindChildTransferByTo(inst1.GetUserTokenInAccount().PublicKey)
	swap.OutputTransfer = in.FindChildTransferByFrom(inst1.GetUserTokenOutAccount().PublicKey)
	in.Event = []interface{}{swap}
	return nil
}
func ParseTransferOwner(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseUnpause(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}
func ParseWithdraw(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*stable_swap.Withdraw)
	removeLiquidity := &types.RemoveLiquidity{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPoolAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	removeLiquidity.TokenATransfer = in.FindChildTransferByFrom(inst1.GetVaultTokenAAccount().PublicKey)
	in.Event = []interface{}{removeLiquidity}
	// log.Logger.Info("ignore parse withdraw", "program", stable_swap.ProgramName)
	return nil
}

// Default
func ParseDefault(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

// Fault
func ParseFault(inst *stable_swap.Instruction, in *types.Instruction, meta *types.Meta) error {
	panic("not supported")
}
