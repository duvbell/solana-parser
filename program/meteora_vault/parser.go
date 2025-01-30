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
	RegisterParser(uint64(meteora_vault.Instruction_Initialize.Uint32()), ParseInitialize)
	RegisterParser(uint64(meteora_vault.Instruction_EnableVault.Uint32()), ParseEnableVault)
	RegisterParser(uint64(meteora_vault.Instruction_SetOperator.Uint32()), ParseSetOperator)
	RegisterParser(uint64(meteora_vault.Instruction_InitializeStrategy.Uint32()), ParseInitializeStrategy)
	RegisterParser(uint64(meteora_vault.Instruction_RemoveStrategy.Uint32()), ParseRemoveStrategy)
	RegisterParser(uint64(meteora_vault.Instruction_RemoveStrategy2.Uint32()), ParseRemoveStrategy2)
	RegisterParser(uint64(meteora_vault.Instruction_CollectDust.Uint32()), ParseCollectDust)
	RegisterParser(uint64(meteora_vault.Instruction_AddStrategy.Uint32()), ParseAddStrategy)
	RegisterParser(uint64(meteora_vault.Instruction_DepositStrategy.Uint32()), ParseDepositStrategy)
	RegisterParser(uint64(meteora_vault.Instruction_WithdrawStrategy.Uint32()), ParseWithdrawStrategy)
	RegisterParser(uint64(meteora_vault.Instruction_Withdraw2.Uint32()), ParseWithdraw2)
	RegisterParser(uint64(meteora_vault.Instruction_Deposit.Uint32()), ParseDeposit)
	RegisterParser(uint64(meteora_vault.Instruction_Withdraw.Uint32()), ParseWithdraw)
	RegisterParser(uint64(meteora_vault.Instruction_WithdrawDirectlyFromStrategy.Uint32()), ParseWithdrawDirectlyFromStrategy)
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

func ParseInitialize(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseEnableVault(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSetOperator(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitializeStrategy(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseRemoveStrategy(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseRemoveStrategy2(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCollectDust(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseAddStrategy(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseDepositStrategy(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseWithdrawStrategy(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseWithdraw2(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	in.Event = in.Children[0].Event
}
func ParseDeposit(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	in.Event = in.Children[0].Event
}
func ParseWithdraw(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	in.Event = in.Children[0].Event
}
func ParseWithdrawDirectlyFromStrategy(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Default
func ParseDefault(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *meteora_vault.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
