package stable_vault

import (
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
	program.RegisterParser(stable_vault.ProgramID, ProgramParser)
	RegisterParser(uint64(stable_vault.Instruction_Withdraw.Uint32()), ParseWithdraw)
	RegisterParser(uint64(stable_vault.Instruction_WithdrawV2.Uint32()), ParseWithdrawV2)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := stable_vault.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

// Withdraw
func ParseWithdraw(inst *stable_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*stable_vault.Withdraw)
	if *inst1.BeneficiaryAmount == 0 {
		// no fee
		in.Event = in.Children[0].Event
	} else {
		in.Event = in.Children[1].Event
	}
}

// WithdrawV2
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
