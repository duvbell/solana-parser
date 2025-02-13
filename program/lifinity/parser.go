package lifinity

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/lifinity_v2"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *lifinity_v2.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

var (
	Instruction_UpdateTargetPriceBufferParam = ag_binary.TypeID([8]byte{175, 241, 237, 224, 206, 9, 187, 30})
	Instruction_UpdateConfigSpreadParam      = ag_binary.TypeID([8]byte{162, 176, 79, 86, 208, 217, 104, 246})
)

func init() {
	program.RegisterParser(lifinity_v2.ProgramID, lifinity_v2.ProgramName, program.Swap, ProgramParser)
	RegisterParser(uint64(lifinity_v2.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(lifinity_v2.Instruction_DepositAllTokenTypes.Uint32()), ParseDepositAllTokenTypes)
	RegisterParser(uint64(lifinity_v2.Instruction_WithdrawAllTokenTypes.Uint32()), ParseWithdrawAllTokenTypes)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	dec := ag_binary.NewBorshDecoder(in.Instruction.Data)
	typeID, err := dec.ReadTypeID()
	if typeID == Instruction_UpdateTargetPriceBufferParam || typeID == Instruction_UpdateConfigSpreadParam {
		return nil
	}
	inst, err := lifinity_v2.DecodeInstruction(in.AccountMetas(meta.Accounts), in.Instruction.Data)
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

// Swap
func ParseSwap(inst *lifinity_v2.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*lifinity_v2.Swap)
	swap := &types.Swap{
		Dex:  in.Instruction.ProgramId,
		Pool: inst1.GetAmmAccount().PublicKey,
		User: inst1.GetAuthorityAccount().PublicKey,
	}
	if *inst1.AmountIn > 0 {
		// the first one is user deposit
		// the second is vault withdraw
		transfers := in.FindChildrenTransfers()
		swap.InputTransfer = transfers[0]
		swap.OutputTransfer = transfers[1]
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
