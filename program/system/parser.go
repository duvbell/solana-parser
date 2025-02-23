package system

import (
	"encoding/binary"
	"errors"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

var (
	programId = solana.SystemProgramID
	Parsers   = make(map[uint64]Parser, 0)
)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

type Parser func(in *system.Instruction, transaction *types.Transaction, index int) error

func init() {
	program.RegisterParser(programId, "system", program.Token, 0, ProgramParser)
	RegisterParser(uint64(system.Instruction_Transfer), ParseTransfer)
}

func ProgramParser(transaction *types.Transaction, index int) error {
	in := transaction.Instructions[index]
	dec := ag_binary.NewBorshDecoder(in.Raw.DataBytes)
	typeID, err := dec.ReadUint32(binary.LittleEndian)
	if _, ok := Parsers[uint64(typeID)]; !ok {
		return nil
	}
	inst, err := system.DecodeInstruction(in.Raw.AccountValues, in.Raw.DataBytes)
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

func ParseTransfer(inst *system.Instruction, transaction *types.Transaction, index int) error {
	inst1 := inst.Impl.(*system.Transfer)
	in := transaction.Instructions[index]
	transfer := &types.Transfer{
		Mint: solana.MustPublicKeyFromBase58("11111111111111111111111111111111"),
		From: inst1.GetFundingAccount().PublicKey,
		To:   inst1.GetRecipientAccount().PublicKey,
	}
	if inst1.Lamports != nil {
		transfer.Amount = *inst1.Lamports
	}
	in.Event = []interface{}{transfer}
	return nil
}
