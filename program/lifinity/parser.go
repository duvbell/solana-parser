package lifinity

import (
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go/programs/lifinity_v2"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *lifinity_v2.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(lifinity_v2.ProgramID, ProgramParser)
	RegisterParser(uint64(lifinity_v2.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(lifinity_v2.Instruction_DepositAllTokenTypes.Uint32()), ParseDepositAllTokenTypes)
	RegisterParser(uint64(lifinity_v2.Instruction_WithdrawAllTokenTypes.Uint32()), ParseWithdrawAllTokenTypes)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) {
	inst, err := lifinity_v2.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

// Swap
func ParseSwap(inst *lifinity_v2.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*lifinity_v2.Swap)
	// the first one is user deposit
	// the second is vault withdraw
	transfers := in.FindChildrenTransfers()
	swap := &types.Swap{
		Pool:           inst1.GetAmmAccount().PublicKey,
		User:           inst1.GetAuthorityAccount().PublicKey,
		TokenATransfer: transfers[0],
		TokenBTransfer: transfers[1],
	}
	in.Event = []interface{}{swap}
}

func ParseDepositAllTokenTypes(inst *lifinity_v2.Instruction, in *types.Instruction, meta *types.Meta) {
	// add liquidity
	log.Logger.Info("ignore parse deposit all token types", "program", lifinity_v2.ProgramName)
}

func ParseWithdrawAllTokenTypes(inst *lifinity_v2.Instruction, in *types.Instruction, meta *types.Meta) {
	// remove liquidity
	log.Logger.Info("ignore parse withdraw all token types", "program", lifinity_v2.ProgramName)
}

// Default
func ParseDefault(inst *lifinity_v2.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *lifinity_v2.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
