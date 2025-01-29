package meteora_vault

import (
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/meteora_vault"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(meteora_vault.ProgramID, ProgramParser)
	RegisterParser(uint64(meteora_vault.Instruction_Withdraw.Uint32()), ParseWithdraw)
	RegisterParser(uint64(meteora_vault.Instruction_Withdraw2.Uint32()), ParseWithdraw2)
	RegisterParser(uint64(meteora_vault.Instruction_Deposit.Uint32()), ParseDeposit)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := meteora_vault.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

// Deposit
func ParseDeposit(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	in.Event = in.Children[0].Event
}

// Withdraw
func ParseWithdraw(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	in.Event = in.Children[0].Event
}

// Withdraw2
func ParseWithdraw2(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	in.Event = in.Children[0].Event
}

// Default
func ParseDefault(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
