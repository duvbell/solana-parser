package solfi

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/solfi"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *solfi.Instruction, in *types.Instruction, meta *types.Meta) error

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(solfi.ProgramID, solfi.ProgramName, program.Swap, 1, ProgramParser)
	RegisterParser(uint64(solfi.Instruction_CreatePair), ParseDefault)
	RegisterParser(uint64(solfi.Instruction_CreatePairV2), ParseDefault)
	RegisterParser(uint64(solfi.Instruction_UpdateConcentration), ParseDefault)
	RegisterParser(uint64(solfi.Instruction_UpdateVersion), ParseDefault)
	RegisterParser(uint64(solfi.Instruction_UpdateFeeParams), ParseDefault)
	RegisterParser(uint64(solfi.Instruction_UpdateOracles), ParseDefault)
	RegisterParser(uint64(solfi.Instruction_WithdrawFees), ParseDefault)
	RegisterParser(uint64(solfi.Instruction_Swap), ParseSwap)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	dec := ag_binary.NewBorshDecoder(in.RawInstruction.DataBytes)
	typeID, err := dec.ReadUint8()
	if typeID != solfi.Instruction_Swap {
		return nil
	}
	inst, err := solfi.DecodeInstruction(in.RawInstruction.AccountValues, in.RawInstruction.DataBytes)
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

// Swap
func ParseSwap(inst *solfi.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*solfi.Swap)
	swap := &types.Swap{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetPairAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	if *inst1.A2B == 0 {
		swap.InputTransfer = in.FindChildTransferByTo(inst1.GetPoolTokenAccountAAccount().PublicKey)
		swap.OutputTransfer = in.FindChildTransferByFrom(inst1.GetPoolTokenAccountBAccount().PublicKey)
	} else {
		swap.InputTransfer = in.FindChildTransferByTo(inst1.GetPoolTokenAccountAAccount().PublicKey)
		swap.OutputTransfer = in.FindChildTransferByFrom(inst1.GetPoolTokenAccountBAccount().PublicKey)
	}
	in.Event = []interface{}{swap}
	return nil
}

// Default
func ParseDefault(inst *solfi.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

// Fault
func ParseFault(inst *solfi.Instruction, in *types.Instruction, meta *types.Meta) error {
	panic("not supported")
}
