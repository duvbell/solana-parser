package raydium_cp

import (
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
	panic("not supported")
}
func ParseUpdateAmmConfig(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseUpdatePoolStatus(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCollectProtocolFee(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseCollectFundFee(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseInitialize(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseDeposit(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseWithdraw(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSwapBaseInput(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
func ParseSwapBaseOutput(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}

// Default
func ParseDefault(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *raydium_cp.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
