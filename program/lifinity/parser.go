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

type Parser func(inst *lifinity_v2.Instruction, transaction *types.Transaction, index int) error

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

var (
	Instruction_UpdateTargetPriceBufferParam = ag_binary.TypeID([8]byte{175, 241, 237, 224, 206, 9, 187, 30})
	Instruction_UpdateConfigSpreadParam      = ag_binary.TypeID([8]byte{162, 176, 79, 86, 208, 217, 104, 246})
)

func init() {
	program.RegisterParser(lifinity_v2.ProgramID, lifinity_v2.ProgramName, program.Swap, 1, ProgramParser)
	RegisterParser(uint64(lifinity_v2.Instruction_Swap.Uint32()), ParseSwap)
	RegisterParser(uint64(lifinity_v2.Instruction_DepositAllTokenTypes.Uint32()), ParseDepositAllTokenTypes)
	RegisterParser(uint64(lifinity_v2.Instruction_WithdrawAllTokenTypes.Uint32()), ParseWithdrawAllTokenTypes)
}

func ProgramParser(transaction *types.Transaction, index int) error {
	in := transaction.Instructions[index]
	dec := ag_binary.NewBorshDecoder(in.Raw.DataBytes)
	typeID, err := dec.ReadTypeID()
	if typeID == Instruction_UpdateTargetPriceBufferParam || typeID == Instruction_UpdateConfigSpreadParam {
		return nil
	}
	inst, err := lifinity_v2.DecodeInstruction(in.Raw.AccountValues, in.Raw.DataBytes)
	if err != nil {
		return err
	}
	id := uint64(inst.TypeID.Uint32())
	parser, ok := Parsers[id]
	if !ok {
		return errors.New("parser not found")
	}
	return parser(inst, transaction, index)
}

// Swap
func ParseSwap(inst *lifinity_v2.Instruction, transaction *types.Transaction, index int) error {
	inst1 := inst.Impl.(*lifinity_v2.Swap)
	in := transaction.Instructions[index]
	swap := &types.Swap{
		Dex:  in.Raw.ProgID,
		Pool: inst1.GetAmmAccount().PublicKey,
		User: inst1.GetAuthorityAccount().PublicKey,
	}
	in.Event = []interface{}{swap}
	log.Logger.Info("ignore parse swap", "program", lifinity_v2.ProgramName)
	return nil
}

func ParseDepositAllTokenTypes(inst *lifinity_v2.Instruction, transaction *types.Transaction, index int) error {
	// add liquidity
	log.Logger.Info("ignore parse deposit all token types", "program", lifinity_v2.ProgramName)
	return nil
}

func ParseWithdrawAllTokenTypes(inst *lifinity_v2.Instruction, transaction *types.Transaction, index int) error {
	// remove liquidity
	log.Logger.Info("ignore parse withdraw all token types", "program", lifinity_v2.ProgramName)
	return nil
}

// Default
func ParseDefault(inst *lifinity_v2.Instruction, transaction *types.Transaction, index int) error {
	return nil
}

// Fault
func ParseFault(inst *lifinity_v2.Instruction, transaction *types.Transaction, index int) error {
	panic("not supported")
}
