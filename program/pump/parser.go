package pump

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/pumpfun"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

var (
	Instruction_AnchorSelfCPILog = ag_binary.TypeID([8]byte{228, 69, 165, 46, 81, 203, 154, 29})
)

func init() {
	program.RegisterParser(pumpfun.ProgramID, ProgramParser)
	RegisterParser(uint64(pumpfun.Instruction_Initialize.Uint32()), ParseInitialize)
	RegisterParser(uint64(pumpfun.Instruction_Create.Uint32()), ParseCreate)
	RegisterParser(uint64(pumpfun.Instruction_Buy.Uint32()), ParseBuy)
	RegisterParser(uint64(pumpfun.Instruction_Sell.Uint32()), ParseSell)
	RegisterParser(uint64(pumpfun.Instruction_Withdraw.Uint32()), ParseWithdraw)
	RegisterParser(uint64(pumpfun.Instruction_SetParams.Uint32()), ParseDefault)
	RegisterParser(uint64(pumpfun.Instruction_AnchorSelfCPILog.Uint32()), ParseDefault)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	dec := ag_binary.NewBorshDecoder(in.Instruction.Data)
	typeID, err := dec.ReadTypeID()
	if typeID == Instruction_AnchorSelfCPILog {
		return nil
	}
	inst, err := pumpfun.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
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

// Initialize
func ParseInitialize(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize", "program", pumpfun.ProgramName)
}

// Create
func ParseCreate(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse create", "program", pumpfun.ProgramName)
}

// Buy
func ParseBuy(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse buy", "program", pumpfun.ProgramName)
}

// Sell
func ParseSell(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse sell", "program", pumpfun.ProgramName)
}

// Sell
func ParseWithdraw(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse withdraw", "program", pumpfun.ProgramName)
}

// Default
func ParseDefault(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *pumpfun.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
