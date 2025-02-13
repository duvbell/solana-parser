package stable_vault

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/stable_vault"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(stable_vault.ProgramID, stable_vault.ProgramName, program.StableSwap, ProgramParser)
	RegisterParser(uint64(stable_vault.Instruction_AcceptAdmin.Uint32()), ParseAcceptAdmin)
	RegisterParser(uint64(stable_vault.Instruction_ChangeBeneficiary.Uint32()), ParseChangeBeneficiary)
	RegisterParser(uint64(stable_vault.Instruction_ChangeBeneficiaryFee.Uint32()), ParseChangeBeneficiaryFee)
	RegisterParser(uint64(stable_vault.Instruction_Initialize.Uint32()), ParseInitialize)
	RegisterParser(uint64(stable_vault.Instruction_Pause.Uint32()), ParsePause)
	RegisterParser(uint64(stable_vault.Instruction_RejectAdmin.Uint32()), ParseRejectAdmin)
	RegisterParser(uint64(stable_vault.Instruction_TransferAdmin.Uint32()), ParseTransferAdmin)
	RegisterParser(uint64(stable_vault.Instruction_Unpause.Uint32()), ParseUnpause)
	RegisterParser(uint64(stable_vault.Instruction_Withdraw.Uint32()), ParseWithdraw)
	RegisterParser(uint64(stable_vault.Instruction_WithdrawV2.Uint32()), ParseWithdrawV2)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	inst, err := stable_vault.DecodeInstruction(in.AccountMetas(meta.Accounts), in.Instruction.Data)
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

func ParseAcceptAdmin(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseChangeBeneficiary(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseChangeBeneficiaryFee(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseInitialize(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize", "program", stable_vault.ProgramName)
}
func ParsePause(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseRejectAdmin(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseTransferAdmin(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUnpause(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseWithdraw(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*stable_vault.Withdraw)
	if *inst1.BeneficiaryAmount == 0 {
		// no fee
		in.Event = in.Children[0].Event
	} else {
		in.Event = in.Children[1].Event
	}
}
func ParseWithdrawV2(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*stable_vault.WithdrawV2)
	if *inst1.BeneficiaryAmount == 0 {
		// no fee
		in.Event = in.Children[0].Event
	} else {
		in.Event = in.Children[1].Event
	}
}

// Default
func ParseDefault(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
