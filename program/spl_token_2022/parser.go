package spl_token_2022

import (
	"errors"

	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
)

var (
	programId = solana.Token2022ProgramID
	Parsers   = make(map[uint64]Parser, 0)
)

type Parser func(inst *token.Instruction, in *types.Instruction, meta *types.Meta) error

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(programId, "token2022", program.Token, 0, ProgramParser)
	RegisterParser(uint64(token.Instruction_Transfer), ParseTransfer)
	RegisterParser(uint64(token.Instruction_TransferChecked), ParseTransferChecked)
	RegisterParser(uint64(token.Instruction_MintTo), ParseMint)
	RegisterParser(uint64(token.Instruction_Burn), ParseBurn)
	RegisterParser(uint64(token.Instruction_InitializeAccount), ParseInitializeAccount)
	RegisterParser(uint64(token.Instruction_InitializeAccount3), ParseInitializeAccount3)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	dec := ag_binary.NewBorshDecoder(in.RawInstruction.DataBytes)
	typeID, err := dec.ReadUint8()
	if _, ok := Parsers[uint64(typeID)]; !ok {
		return nil
	}
	inst, err := token.DecodeInstruction(in.RawInstruction.AccountValues, in.RawInstruction.DataBytes)
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

func ParseTransfer(inst *token.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*token.Transfer)
	transfer := &types.Transfer{
		From: inst1.GetSourceAccount().PublicKey,
		To:   inst1.GetDestinationAccount().PublicKey,
	}
	mint := meta.TokenAccounts[transfer.From]
	transfer.Mint = mint.Mint
	if inst1.Amount != nil {
		transfer.Amount = *inst1.Amount
	}
	in.Event = []interface{}{transfer}
	return nil
}

func ParseTransferChecked(inst *token.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*token.TransferChecked)
	transfer := &types.Transfer{
		From: inst1.GetSourceAccount().PublicKey,
		To:   inst1.GetDestinationAccount().PublicKey,
	}
	mint := meta.TokenAccounts[transfer.From]
	transfer.Mint = mint.Mint
	if inst1.Amount != nil {
		transfer.Amount = *inst1.Amount
	}
	in.Event = []interface{}{transfer}
	return nil
}

func ParseMint(inst *token.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*token.MintTo)
	mintTo := &types.MintTo{
		Mint:    inst1.GetMintAccount().PublicKey,
		Account: inst1.GetDestinationAccount().PublicKey,
	}
	if inst1.Amount != nil {
		mintTo.Amount = *inst1.Amount
	}
	in.Event = []interface{}{mintTo}
	return nil
}

func ParseBurn(inst *token.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*token.Burn)
	burn := &types.Burn{
		Mint:    inst1.GetMintAccount().PublicKey,
		Account: inst1.GetSourceAccount().PublicKey,
	}
	if inst1.Amount != nil {
		burn.Amount = *inst1.Amount
	}
	in.Event = []interface{}{burn}
	return nil
}

func ParseInitializeAccount(inst *token.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*token.InitializeAccount)
	init := &types.Initialize{
		Mint:    inst1.GetMintAccount().PublicKey,
		Account: inst1.GetAccount().PublicKey,
		Owner:   inst1.GetOwnerAccount().PublicKey,
	}
	// update token owner & mint by spl token instructions
	meta.TokenAccounts[init.Account] = &types.TokenAccount{
		Owner:     &init.Owner,
		ProgramId: &in.RawInstruction.ProgID,
		Mint:      init.Mint,
	}
	in.Event = []interface{}{init}
	return nil
}

func ParseInitializeAccount3(inst *token.Instruction, in *types.Instruction, meta *types.Meta) error {
	inst1 := inst.Impl.(*token.InitializeAccount3)
	init := &types.Initialize{
		Mint:    inst1.GetMintAccount().PublicKey,
		Account: inst1.GetAccount().PublicKey,
	}
	if inst1.Owner != nil {
		init.Owner = *inst1.Owner
	}
	// update token owner & mint by spl token instructions
	meta.TokenAccounts[init.Account] = &types.TokenAccount{
		Owner:     &init.Owner,
		ProgramId: &in.RawInstruction.ProgID,
		Mint:      init.Mint,
	}
	in.Event = []interface{}{init}
	return nil
}
