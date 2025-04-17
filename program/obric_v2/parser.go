package obricv2

import (
	"errors"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/obric_v2"
	"github.com/solana-parser/types"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *obric_v2.Instruction, in *types.Instruction, meta *types.Meta) error

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

var (
	Instruction_UpdateTargetPriceBufferParam = ag_binary.TypeID([8]byte{175, 241, 237, 224, 206, 9, 187, 30})
	Instruction_UpdateConfigSpreadParam      = ag_binary.TypeID([8]byte{162, 176, 79, 86, 208, 217, 104, 246})
)

func init() {
	//program.RegisterParser(obric_v2.ProgramID, obric_v2.ProgramName, program.Swap, 1, ProgramParser)
	RegisterParser(uint64(obric_v2.Instruction_CreatePair.Uint32()), ParseDefault)
	RegisterParser(uint64(obric_v2.Instruction_CreatePairV2.Uint32()), ParseDefault)
	RegisterParser(uint64(obric_v2.Instruction_UpdateConcentration.Uint32()), ParseDefault)
	RegisterParser(uint64(obric_v2.Instruction_UpdateVersion.Uint32()), ParseDefault)
	RegisterParser(uint64(obric_v2.Instruction_UpdateFeeParams.Uint32()), ParseDefault)
	RegisterParser(uint64(obric_v2.Instruction_UpdateOracles.Uint32()), ParseDefault)
	RegisterParser(uint64(obric_v2.Instruction_WithdrawFees.Uint32()), ParseDefault)
	RegisterParser(uint64(obric_v2.Instruction_Deposit.Uint32()), ParseDefault)
	RegisterParser(uint64(obric_v2.Instruction_Withdraw.Uint32()), ParseDefault)
	RegisterParser(uint64(obric_v2.Instruction_SwapXToY.Uint32()), ParseSwapXToY)
	RegisterParser(uint64(obric_v2.Instruction_SwapYToX.Uint32()), ParseSwapYToX)
	RegisterParser(uint64(obric_v2.Instruction_Swap.Uint32()), ParseSwap)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	dec := ag_binary.NewBorshDecoder(in.RawInstruction.DataBytes)
	typeID, err := dec.ReadTypeID()
	if typeID == Instruction_UpdateTargetPriceBufferParam || typeID == Instruction_UpdateConfigSpreadParam {
		return nil
	}
	inst, err := obric_v2.DecodeInstruction(in.RawInstruction.AccountValues, in.RawInstruction.DataBytes)
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
func ParseSwap(inst *obric_v2.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*obric_v2.Swap)
	swap := &types.Swap{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetTradingPairAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	if *inst1.IsXToY {
		swap.InputTransfer = in.FindNextTransferByTo(inst1.GetReserveXAccount().PublicKey)
		swap.OutputTransfer = in.FindNextTransferByFrom(inst1.GetMintYAccount().PublicKey)
	} else {
		swap.InputTransfer = in.FindNextTransferByTo(inst1.GetReserveYAccount().PublicKey)
		swap.OutputTransfer = in.FindNextTransferByFrom(inst1.GetMintXAccount().PublicKey)
	}
	in.Event = []interface{}{swap}
	return nil
}

func ParseSwapXToY(inst *obric_v2.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*obric_v2.SwapXToY)
	swap := &types.Swap{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetTradingPairAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	swap.InputTransfer = in.FindNextTransferByTo(inst1.GetReserveXAccount().PublicKey)
	swap.OutputTransfer = in.FindNextTransferByFrom(inst1.GetMintYAccount().PublicKey)
	in.Event = []interface{}{swap}
	return nil
}

func ParseSwapYToX(inst *obric_v2.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*obric_v2.SwapYToX)
	swap := &types.Swap{
		Dex:  in.RawInstruction.ProgID,
		Pool: inst1.GetTradingPairAccount().PublicKey,
		User: inst1.GetUserAccount().PublicKey,
	}
	swap.InputTransfer = in.FindNextTransferByTo(inst1.GetReserveYAccount().PublicKey)
	swap.OutputTransfer = in.FindNextTransferByFrom(inst1.GetMintXAccount().PublicKey)
	in.Event = []interface{}{swap}
	return nil
}

// Default
func ParseDefault(inst *obric_v2.Instruction, in *types.Instruction, meta *types.Meta) error {
	return nil
}

// Fault
func ParseFault(inst *obric_v2.Instruction, in *types.Instruction, meta *types.Meta) error {
	panic("not supported")
}
